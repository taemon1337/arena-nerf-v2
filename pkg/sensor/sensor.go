package sensor

import (
  "fmt"
  "log"
  "time"
  "strings"
  "context"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

type Sensor struct {
  conf          *config.SensorConfig
  gamechan      *game.GameChannel
  *log.Logger
}

func NewSensor(cfg *config.SensorConfig, gamechan *game.GameChannel, logger *log.Logger) *Sensor {
  return &Sensor{
    conf:     cfg,
    gamechan: gamechan,
    Logger:   log.New(logger.Writer(), "[SENSOR]: ", logger.Flags()),
  }
}

func (s *Sensor) Start(ctx context.Context) error {
  s.Printf("starting sensor")
  for {
    select {
    case evt := <-s.gamechan.SensorChan:
      s.Printf("sensor received game event: %s", evt)
      switch evt.Event {
        case constants.SENSOR_HIT_REQUEST:
          parts := strings.Split(string(evt.Payload), constants.SPLIT)
          if len(parts) != 2 {
            s.Printf("error parsing sensor hit request: %s (should be <1-4>:<num>)", string(evt.Payload))
            continue
          }

          if parts[0] == "1" || parts[0] == "2" || parts[0] == "3" || parts[0] == "4" {
            s.Printf("generating sensor hit by request: %s", evt)
            s.SensorHit(parts[0])
          } else {
            s.Printf("error parsing sensor hit request: %s (should be <1-4>:<num>)", string(evt.Payload))
            continue
          }
        default:
          s.Printf("unrecognized sensor event: %s", evt)
      }
    case <-ctx.Done():
      s.Printf("stopping sensor")
      return ctx.Err()
    default:
      time.Sleep(3 * time.Second) // do something later
    }
  }
}

func (s *Sensor) SensorHit(sensorid string) {
  pay := fmt.Sprintf("%s%s1", sensorid, constants.SPLIT) // sensor hits are always a single hit, 1
  s.gamechan.GameChan <- game.NewGameEvent(constants.SENSOR_HIT, []byte(pay))
}

func (s *Sensor) IsTestSensor() bool {
  return strings.HasPrefix(s.conf.Id, constants.TEST_SENSOR_PREFIX)
}
