package game

import (
  "log"
  "time"
  "context"

  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type GameSimulation struct {
  conf          *config.Config
  gamechan      *GameChannel
  state         string
  scoreboard    map[string]int
  *log.Logger
}

func NewGameSimulation(cfg *config.Config, gamechan *GameChannel, scoreboard map[string]int, logger *log.Logger) *GameSimulation {
  return &GameSimulation{
    conf:           cfg,
    gamechan:       gamechan,
    state:          constants.GAME_STATE_INIT,
    scoreboard:     scoreboard,
    Logger:         log.New(logger.Writer(), "simulation: ", logger.Flags()),
  }
}

func (g *GameSimulation) Start(ctx context.Context) error {
  g.Printf("starting game %s", g)
  g.state = constants.GAME_STATE_RUNNING
  for {
    select {
    case evt := <- g.gamechan.EventChan:
      g.Printf("game simulation received game event: %s", evt.Event)
    case <-ctx.Done():
      g.Printf("stopping game %s", g)
      g.state = constants.GAME_STATE_ENDED
      return ctx.Err()
    case <-time.After(3 * time.Second):
      g.GenerateRandomGameEvent()
    }
  }
}

func (g *GameSimulation) GenerateRandomGameEvent() {
  g.Printf("generating random game event")
  g.gamechan.RequestChan <- NewGameEvent(constants.RANDOM_TEAM_HIT, []byte(""))
}

func (g *GameSimulation) Stop(ctx context.Context) error {
  ctx.Done()
  return ctx.Err()
}

func (g *GameSimulation) Completed() bool {
  return false
}

func (g *GameSimulation) Running() bool {
  return g.state == constants.GAME_STATE_RUNNING
}

func (g *GameSimulation) String() string {
  return constants.GAME_MODE_SIMULATION
}
