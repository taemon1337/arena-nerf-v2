package node

import (
  "log"
  "time"
  "strings"
  "context"

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
  *log.Logger
}

func NewNode(cfg *config.Config, gamechan *game.GameChannel, logger *log.Logger) *Node {
  logger = log.New(logger.Writer(), "[NODE]: ", logger.Flags())

  return &Node{
    conf:     cfg,
    conn:     connector.NewConnector(cfg, logger),
    gamechan: gamechan,
    sensor:   sensor.NewSensor(cfg, gamechan, logger),
    Logger:   logger,
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
      time.Sleep(5 * time.Second)
      return n.conn.Join()
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
