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

// https://github.com/hybridgroup/gobot/blob/release/drivers/gpio/rgb_led_driver.go

type Sensor struct {
  conf          *config.SensorConfig
  gamechan      *game.GameChannel
  Sensorchan    chan game.GameEvent
  led           *SensorLed
  *log.Logger
}

func NewSensor(cfg *config.SensorConfig, gamechan *game.GameChannel, logger *log.Logger) *Sensor {
  return &Sensor{
    conf:         cfg,
    gamechan:     gamechan,
    Sensorchan:   make(chan game.GameEvent, 5),
    led:          NewSensorLed(constants.COLOR_YELLOW),
    Logger:   log.New(logger.Writer(), fmt.Sprintf("[sensor:%s]: ", cfg.Id), logger.Flags()),
  }
}

func (s *Sensor) Start(ctx context.Context) error {
  s.Printf("starting sensor")
  for {
    select {
    case evt := <-s.Sensorchan:
      switch evt.Event {
        case constants.SENSOR_HIT:
          s.Printf("sensor received sensor hit game event: %s", evt)
          s.SensorHit(s.conf.Id)
        case constants.SENSOR_COLOR:
          s.Printf("sensor received sensor color game event: %s", evt)
          s.led.SetColor(string(evt.Payload))
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
  pay := strings.Join([]string{sensorid, s.led.GetColor(), "1"}, constants.SPLIT)
  s.gamechan.GameChan <- game.NewGameEvent(constants.SENSOR_HIT, []byte(pay))
}

func (s *Sensor) IsTestSensor() bool {
  return strings.HasPrefix(s.conf.Id, constants.TEST_SENSOR_PREFIX)
}

func (s *Sensor) Led() *SensorLed {
  return s.led
}
