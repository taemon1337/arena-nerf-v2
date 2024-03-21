package node

import (
  "log"
  "fmt"
  "time"
  "sync"
  "strings"
  "context"
  "math/rand"
  "encoding/json"

  "golang.org/x/sync/errgroup"
  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/arena-nerf/pkg/common"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/connector"
  "github.com/taemon1337/arena-nerf/pkg/game"
  "github.com/taemon1337/arena-nerf/pkg/sensor"
)

type Node struct {
  conf          *config.Config
  conn          *connector.Connector
  sensors       map[string]*sensor.Sensor
  gamechan      *game.GameChannel
  nodestate     *NodeState
  nodelock      *sync.Mutex
  *log.Logger
}

func NewNode(cfg *config.Config, gamechan *game.GameChannel, logger *log.Logger) *Node {
  logger = log.New(logger.Writer(), fmt.Sprintf("[%s]: ", cfg.AgentConf.NodeName), logger.Flags())

  return &Node{
    conf:       cfg,
    conn:       connector.NewConnector(cfg, logger),
    gamechan:   gamechan,
    sensors:    map[string]*sensor.Sensor{},
    nodestate:  NewNodeState(cfg.AgentConf.NodeName),
    nodelock:   &sync.Mutex{},
    Logger:     logger,
  }
}

func (n *Node) Start(ctx context.Context) error {
  n.Printf("starting node")
  g, ctx := errgroup.WithContext(ctx)

  if n.conf.EnableConnector {
    err := n.conn.Connect()
    if err != nil {
      return err
    }

    n.conn.RegisterEventHandler(n)

    g.Go(func () error {
      time.Sleep(3 * time.Second)
      return n.conn.Join(ctx)
    })
  }

  if n.conf.EnableSensors {
    for id, cfg := range n.conf.SensorsConf.Configs {
      sensconf := cfg // local variable to ensure proper (even without closure)
      n.Printf("initializing sensor %s", id)
      err := sensconf.Error()
      if err != nil && err != constants.ERR_TEST_SENSOR {
        return err
      }

      n.sensors[id] = sensor.NewSensor(id, sensconf, n.gamechan, n.Logger, n.conf.EnableLeds, n.conf.EnableHits)
    }
  } else {
    n.Printf("sensors disabled")
  }

  for id, _ := range n.sensors {
    sens := n.GetSensorById(id) // local variable needed to store id inside loop (otherwise it will call the same sensor 2x)
    g.Go(func() error {
      return sens.Start(ctx)
    })
    n.Printf("started sensor %s", id)
  }

  g.Go(func() error {
    for {
      select {
      case e := <-n.gamechan.GameChan:
        switch e.Event {
          case constants.SENSOR_HIT:
            sensorid, sensorcolor, hitcount, err := common.ParseNodeHitPayload(e.Payload)
            if err != nil {
              n.Printf("error parsing node hit payload: %s", err)
              continue
            }

            n.Printf("node received sensor hit: %s", e)
            n.nodestate.AddNodeHit(sensorid, sensorcolor, hitcount)
            n.Printf("node recorded sensor hit: %s", e)
            continue
          default:
            n.Printf("node received game event: %s", e)
        }
      case <-ctx.Done():
        n.Printf("stopping node")
        n.Close()
        return ctx.Err()
      }
    }
  })

  err := g.Wait()
  if err != nil {
    n.Printf("node watch ended with error: %s", err)
  }
  return err
}

func (n *Node) GetSensorById(id string) *sensor.Sensor {
  if _, ok := n.sensors[id]; !ok {
    return nil
  }
  return n.sensors[id]
}

func (n *Node) Close() {
  for id, _ := range n.sensors {
    sens := n.GetSensorById(id)
    sens.Close()
  }
}

// handle event are events from serf (over the network)
func (n *Node) HandleEvent(evt serf.Event) {
  if evt.EventType() == serf.EventUser {
    e := evt.(serf.UserEvent)
    switch e.Name {
      case constants.GAME_MODE:
        n.Printf("set game mode to %s", string(e.Payload))
        n.nodestate.Mode = string(e.Payload)
      case constants.GAME_ACTION_BEGIN:
        n.Printf("start game received")
        n.nodestate.Status = constants.GAME_STATUS_RUNNING
      case constants.GAME_ACTION_END:
        n.Printf("end game received")
        n.nodestate.Status = constants.GAME_STATUS_ENDED
      case constants.GAME_TEAMS:
        n.Printf("set game teams - %s", string(e.Payload))
        n.nodestate.SetTeams(string(e.Payload), n.conf.EnableTeamColors)
      case n.NodeEventName(constants.SENSOR_HIT_REQUEST):
        // sensor hits always come directly from sensors, not through the network
        // so in this case, it is a synthetic hit and not a real one
        n.Printf("synthetic sensor hit: %s", e.Name)
        if n.nodestate.Status != constants.GAME_STATUS_RUNNING {
          n.Printf("game is not active - no hits allowed")
          return
        }

        sensorid, hitcount, err := common.ParseSensorHit(e.Payload)
        if err != nil {
          n.Printf("error parsing sensor hit request: %s (should be <sensor-name>:<hit-count>): %s", string(e.Payload), err)
          return
        }

        if err := n.SendEventToSensor(sensorid, game.NewGameEvent(constants.SENSOR_HIT, []byte(fmt.Sprintf("%d", hitcount)))); err != nil {
          n.Printf("error sending event %s to sensor: %s", e.Name, err)
          return
        }
      case n.NodeEventName(constants.SENSOR_COLOR_REQUEST):
        n.Printf("node received sensor color request: %s", e.Name)
        if n.nodestate.Status != constants.GAME_STATUS_RUNNING {
          n.Printf("game is not active - cannot set random sensor color")
          return
        }

        parts := strings.Split(string(e.Payload), constants.SPLIT)
        if len(parts) != 2 {
          n.Printf("error parsing sensor color request: %s (should be <sensor-name>:<color>)", string(e.Payload))
          return
        }

        sensorid := parts[0]
        color := parts[1]

        if sensorid == constants.RANDOM_SENSOR_ID {
          sensorid = n.RandomSensorId()
        }

        if color == constants.RANDOM_COLOR_ID {
          sens := n.GetSensorById(sensorid)
          if sens == nil {
            n.Printf("no sensor found named %s on this node", sensorid)
            return 
          }

          led := sens.Led()
          if led != nil {
            currentcolor := led.GetColor()
            color = n.RandomColor(currentcolor)
          }
        }

        if err := n.SendEventToSensor(sensorid, game.NewGameEvent(constants.SENSOR_COLOR, []byte(color))); err != nil {
          n.Printf("error sending event %s to sensor: %s", e.Name, err)
          return
        }
      case n.NodeEventName(constants.TEAM_HIT):
        n.Printf("NODE EVENT: %s", e.Name)
        if n.nodestate.Status != constants.GAME_STATUS_RUNNING {
          n.Printf("game is not active - no hits allowed")
          return
        }

        parts := strings.Split(string(e.Payload), constants.SPLIT)
        if len(parts) < 2 {
          log.Printf("cannot parse team hit from %s - should be <team>:<count>", string(e.Payload))
        } else {
          team, hits, err := common.ParseTeamHit(e.Payload)
          if err != nil {
            n.Printf("error parsing team hit event %s: %s", e.Name, err)
            return
          }

          n.nodestate.AddTeamHit(team, hits)
        }
      default:
        n.Printf("unrecognized event - %s", e.Name)
    }
  }
  if evt.EventType() == serf.EventQuery {
    var err error = nil
    q := evt.(*serf.Query)
    switch q.Name {
      case constants.NODE_READY:
        err = q.Respond([]byte(constants.NODE_IS_READY))
      case constants.GAME_MODE:
        err = q.Respond([]byte(n.nodestate.Mode))
      case constants.NODE_SCOREBOARD:
        data, err := json.Marshal(n.nodestate.Hits)
        if err != nil {
          log.Printf("cannot marshal node hits: %s", err)
        } else {
          err = q.Respond(data)
        }
      default:
        n.Printf("unrecognized query - %s", q.Name)
    }

    if err != nil {
      n.Printf("error responding to query %s: %s", q.Name, err)
    }
  }
}

func (n *Node) NodeEventName(action string) string {
  return strings.Join([]string{n.conf.AgentConf.NodeName, action}, constants.SPLIT)
}

func (n *Node) SendEventToSensor(sensorid string, e game.GameEvent) error {
  if !n.conf.EnableSensors {
    return constants.ERR_SENSORS_DISABLED
  }

  if len(n.sensors) < 1 {
    return constants.ERR_NO_SENSORS
  }

  if sensorid == constants.RANDOM_SENSOR_ID {
    sensorid = n.RandomSensorId()
  }

  if _, ok := n.sensors[sensorid]; !ok {
    return constants.ERR_NO_SENSOR_BY_NAME
  }

  select {
    case n.sensors[sensorid].SensorChan <- e:
      n.Printf("sent event to sensor: %s", e)
    default:
      n.Printf("sensor chan is full - discarding event: %s", e)
  }
  return nil
}

func (n *Node) RandomSensorId() string {
  i := rand.Intn(len(n.sensors))
  p := 0

  for id, _ := range n.sensors {
    if p == i {
      return id
    }
    p += 1
  }

  return ""
}

func (n *Node) RandomColor(except_color string) string {
  if len(n.nodestate.Colors) == 0 {
    return ""
  }

  if len(n.nodestate.Colors) == 1 {
    return n.nodestate.Colors[0]
  }

  availableColors := []string{}
  for _, color := range n.nodestate.Colors {
    if color != except_color {
      availableColors = append(availableColors, color)
    }
  }

  return availableColors[rand.Intn(len(availableColors))]
}
