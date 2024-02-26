package node

import (
  "log"
  "time"
  "context"

  "golang.org/x/sync/errgroup"

  "github.com/taemon1337/arena-nerf/pkg/config"
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
  logger = log.New(logger.Writer(), "node: ", logger.Flags())

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
