package controller

import (
  "log"
  "time"
  "context"

  "golang.org/x/sync/errgroup"
  "github.com/hashicorp/serf/serf"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/connector"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

type Controller struct {
  conf          *config.Config
  engine        *game.GameEngine
  gamechan      *game.GameChannel
  conn          *connector.Connector
  *log.Logger
}

func NewController(cfg *config.Config, gamechan *game.GameChannel, logger *log.Logger) *Controller {
  logger = log.New(logger.Writer(), "ctrl: ", logger.Flags())

  return &Controller{
    conf:     cfg,
    engine:   game.NewGameEngine(cfg, gamechan, logger),
    gamechan: gamechan,
    conn:     connector.NewConnector(cfg, logger),
    Logger:   logger,
  }
}

func (ctrl *Controller) Start(ctx context.Context) error {
  ctrl.Printf("starting controller")
  g, ctx := errgroup.WithContext(ctx)

  if ctrl.conf.EnableConnector {
    err := ctrl.conn.Connect()
    if err != nil {
      return err
    }

    ctrl.conn.RegisterEventHandler(ctrl)

    g.Go(func () error {
      time.Sleep(5 * time.Second)
      return ctrl.conn.Join()
    })
  }

  if ctrl.conf.EnableGameEngine {
    g.Go(func() error {
      return ctrl.engine.Start(ctx)
    })
  } else {
    ctrl.Printf("game engine disabled")
  }

  g.Go(func() error {
    for {
      select {
      case <-ctx.Done():
        ctrl.Printf("stopping controller")
        return ctx.Err()
      default:
        time.Sleep(3 * time.Second) // do something later
      }
    }
  })

  return g.Wait()
}

func (ctrl *Controller) NewGame(mode string) error {
  return ctrl.engine.NewGame(mode)
}

func (ctrl *Controller) StartGame(ctx context.Context) error {
  return ctrl.engine.StartGame(ctx)
}

func (ctrl *Controller) HandleEvent(e serf.Event) {
  if e.EventType() == serf.EventUser {
    log.Printf("EVENT: %s", e)
  }
  if e.EventType() == serf.EventQuery {
    log.Printf("QUERY: %s", e)
  }
}
