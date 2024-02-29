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
  logger = log.New(logger.Writer(), "[CTRL]: ", logger.Flags())

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
      time.Sleep(3 * time.Second)
      return ctrl.conn.Join(ctx)
    })

    g.Go(func () error {
      return ctrl.ListenToGame(ctx)
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
        ctrl.conn.Shutdown()
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

func (ctrl *Controller) ListenToGame(ctx context.Context) error {
  for {
    select {
    case <-ctx.Done():
      return ctx.Err()
    case e := <-ctrl.gamechan.NodeChan:
      switch e.Event {
        default:
          ctrl.Printf("sending event out to all nodes: %s", e.Event)
          if err := ctrl.conn.UserEvent(e.Event, e.Payload, ctrl.conf.Coalesce); err != nil {
            ctrl.Printf("error sending %s event: %s", e.Event, err)
          }
      }
    case q := <-ctrl.gamechan.QueryChan:
      ctrl.Printf("controller received game query: %s", q)
      switch q.Query {
        default:
          // by default send all queries from game engine to all nodes
          data := map[string][]byte{}
          resp, err := ctrl.conn.Query(q.Query, q.Payload, &serf.QueryParam{FilterTags: q.Tags})
          if err != nil {
            q.Response <- game.NewGameQueryResponse(data, err)
          }
          for r := range resp.ResponseCh() {
            data[r.From] = r.Payload
          }
          q.Response <- game.NewGameQueryResponse(data, err)
      }
    }
  }
}
