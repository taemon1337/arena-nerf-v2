package game

import (
  "log"
  "context"
  "github.com/google/uuid"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type Game interface {
  Id() string
  Mode() string
  Start(ctx context.Context) error
  Stop(ctx context.Context) error
  String() string
}

func NewGame(mode string, cfg *config.Config, gc *GameChannel, logger *log.Logger) Game {
  id := uuid.New().String()
  switch mode {
    case constants.GAME_MODE_SIMULATION:
      return NewGameSimulation(id, mode, cfg, gc, map[string]int{}, logger)
    default:
      logger.Printf("Unsupported game mode %s", mode)
      return nil
  }
}

