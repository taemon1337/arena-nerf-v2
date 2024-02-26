package game

import (
  "log"
  "context"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type Game interface {
  Start(ctx context.Context) error
  Stop(ctx context.Context) error
  Running() bool
  Completed() bool
  String() string
}

func NewGame(mode string, cfg *config.Config, gc *GameChannel, scoreboard map[string]int, logger *log.Logger) Game {
  switch mode {
    case constants.GAME_MODE_SIMULATION:
      return NewGameSimulation(cfg, gc, scoreboard, logger)
    default:
      logger.Printf("Unsupported game mode %s", mode)
      return nil
  }
}

