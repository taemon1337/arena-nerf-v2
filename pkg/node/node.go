package node

import (
  "log"
  "fmt"
  "time"
  "strconv"
  "strings"
  "context"
  "encoding/json"

  "golang.org/x/sync/errgroup"
  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/connector"
  "github.com/taemon1337/arena-nerf/pkg/game"
  "github.com/taemon1337/arena-nerf/pkg/sensor"
)

type Node struct {
  conf          *config.Config
  conn          *connector.Connector
  sensor        *sensor.Sensor
  gamechan      *game.GameChannel
  nodestate     *NodeState
  *log.Logger
}

func NewNode(cfg *config.Config, gamechan *game.GameChannel, logger *log.Logger) *Node {
  logger = log.New(logger.Writer(), fmt.Sprintf("[%s]: ", cfg.AgentConf.NodeName), logger.Flags())

  return &Node{
    conf:       cfg,
    conn:       connector.NewConnector(cfg, logger),
    gamechan:   gamechan,
    sensor:     sensor.NewSensor(cfg, gamechan, logger),
    nodestate:  NewNodeState(cfg.AgentConf.NodeName),
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

  if n.conf.EnableSensor {
    g.Go(func() error {
      return n.sensor.Start(ctx)
    })
  } else {
    n.Printf("sensor disabled")
  }

  g.Go(func() error {
    for {
      select {
      case <-ctx.Done():
        n.Printf("stopping node")
        return ctx.Err()
      default:
        time.Sleep(3 * time.Second) // do something later
      }
    }
  })

  return g.Wait()
}

func (n *Node) HandleEvent(evt serf.Event) {
  if evt.EventType() == serf.EventUser {
    e := evt.(serf.UserEvent)
    switch e.Name {
      case constants.GAME_MODE:
        n.Printf("set game mode to %s", string(e.Payload))
        n.nodestate.SetMode(string(e.Payload))
      case constants.GAME_ACTION_BEGIN:
        n.Printf("start game received")
        n.nodestate.SetStatus(constants.GAME_STATUS_RUNNING)
      case constants.GAME_ACTION_END:
        n.Printf("end game received")
        n.nodestate.SetStatus(constants.GAME_STATUS_ENDED)
      case constants.GAME_TEAMS:
        n.Printf("set game teams - %s", string(e.Payload))
        n.nodestate.SetTeams(string(e.Payload))
      case n.NodeEventName(constants.TEAM_HIT):
        n.Printf("NODE EVENT: %s", e.Name)
        if n.nodestate.Status() != constants.GAME_STATUS_RUNNING {
          n.Printf("game is not active - no hits allowed")
          return
        }

        parts := strings.Split(string(e.Payload), constants.SPLIT)
        if len(parts) < 2 {
          log.Printf("cannot parse team hit from %s - should be <team>:<count>", string(e.Payload))
        } else {
          hits, err := strconv.Atoi(parts[1])
          if err != nil {
            log.Printf("cannot parse team hit from %s - %s", string(e.Payload), err)
          } else {
            n.nodestate.AddTeamHit(parts[0], hits)
            n.nodestate.AddNodeHit(hits)
/*
            if n.HasSensor() {
              n.sensor.NodeTeamHit(constants.TEAM_HIT, e.Payload)
            }
*/
          }
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
        err = q.Respond([]byte(n.nodestate.GetMode()))
      case constants.NODE_SCOREBOARD:
        data, err := json.Marshal(n.nodestate.Hits())
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
