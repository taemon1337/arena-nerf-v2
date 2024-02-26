package main

import (
  "os"
  "log"
  "context"
  "syscall"
  "os/signal"

  "golang.org/x/sync/errgroup"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/game"
  "github.com/taemon1337/arena-nerf/pkg/controller"
  "github.com/taemon1337/arena-nerf/pkg/node"
)

func main() {
  logger := log.New(os.Stdout, "main: ", log.LstdFlags)
  cfg := config.NewConfig(logger)
  gamechan := game.NewGameChannel()

  if err := cfg.Flags(); err != nil {
    logger.Fatal(err)
  }

  logger.Printf(cfg.String())

  ctx, cancel := context.WithCancel(context.Background())
  defer cancel()

  g, ctx := errgroup.WithContext(ctx)

  sigCh := make(chan os.Signal, 2)
  signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(sigCh)

  // watch for signals
  g.Go(func() error {
    select {
    case <-sigCh:
      logger.Printf("signal received, stopping...")
      cancel()
    }
    return nil
  })

  // if controller enabled
  if cfg.EnableController {
    ctrl := controller.NewController(cfg, gamechan, logger)

    g.Go(func() error {
      return ctrl.Start(ctx)
    })

    if cfg.EnableSimulation {
      g.Go(func() error {
        ctrl.NewGame(constants.GAME_MODE_SIMULATION)
        return ctrl.StartGame(ctx)
      })
    }
  }

  // if node enabled
  if cfg.EnableNode {
    nod := node.NewNode(cfg, gamechan, logger)

    g.Go(func() error {
      return nod.Start(ctx)
    })
  }

  // wait for all threads to return, logging first non-nil error
  if err := g.Wait(); err != nil {
    logger.Fatal(err)
  } else {
    logger.Printf("Shutdown complete without error.")
  }
}
