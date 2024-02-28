package game

import (
  "log"
  "time"
  "context"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type GameSimulation struct {
  id            string
  mode          string
  conf          *config.Config
  gamechan      *GameChannel
  scoreboard    map[string]int
  *log.Logger
}

func NewGameSimulation(id, mode string, cfg *config.Config, gamechan *GameChannel, scoreboard map[string]int, logger *log.Logger) *GameSimulation {
  return &GameSimulation{
    id:             id,
    mode:           mode,
    conf:           cfg,
    gamechan:       gamechan,
    scoreboard:     scoreboard,
    Logger:         log.New(logger.Writer(), "[SIMULATION]: ", logger.Flags()),
  }
}

func (g *GameSimulation) Id() string {
  return g.id
}

func (g *GameSimulation) Mode() string {
  return g.mode
}

func (g *GameSimulation) Start(ctx context.Context) error {
  g.Printf("starting game %s", g)
  for {
    select {
    case evt := <- g.gamechan.EventChan:
      g.Printf("game simulation received game event: %s", evt.Event)
    case <-ctx.Done():
      g.gamechan.RequestChan <- NewGameEvent(constants.GAME_ACTION_END, []byte("stopping game - context done"))
      return ctx.Err()
    case <-time.After(3 * time.Second):
      g.gamechan.RequestChan <- NewGameEvent(constants.RANDOM_TEAM_HIT, []byte("simulating random game event"))
    }
  }
}

func (g *GameSimulation) Stop(ctx context.Context) error {
  ctx.Done()
  return ctx.Err()
}

func (g *GameSimulation) String() string {
  return constants.GAME_MODE_SIMULATION
}
