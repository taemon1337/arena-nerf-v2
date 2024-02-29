package sensor

import (
  "log"
  "time"
  "context"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

type Sensor struct {
  conf          *config.Config
  gc            *game.GameChannel
  *log.Logger
}

func NewSensor(cfg *config.Config, gamechan *game.GameChannel, logger *log.Logger) *Sensor {
  return &Sensor{
    conf:     cfg,
    gc:       gamechan,
    Logger:   log.New(logger.Writer(), "[SENSOR]: ", logger.Flags()),
  }
}

func (s *Sensor) Start(ctx context.Context) error {
  s.Printf("starting sensor")
  for {
    select {
    case evt := <-s.gc.SensorChan:
      s.Printf("sensor received game event: %s", evt)
    case <-ctx.Done():
      s.Printf("stopping sensor")
      return ctx.Err()
    default:
      time.Sleep(3 * time.Second) // do something later
    }
  }
}
