package sensor

import (
  "fmt"
  "log"
  "time"
  "strings"
  "context"
  "golang.org/x/sync/errgroup"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

// https://github.com/hybridgroup/gobot/blob/release/drivers/gpio/rgb_led_driver.go

type Sensor struct {
  conf          *config.SensorConfig
  gamechan      *game.GameChannel
  SensorChan    chan game.GameEvent
  led           *SensorLed
  hit           *SensorHitInput
  *log.Logger
}

func NewSensor(cfg *config.SensorConfig, gamechan *game.GameChannel, logger *log.Logger) *Sensor {
  sensorchan := make(chan game.GameEvent, constants.CHANNEL_WIDTH)
  return &Sensor{
    conf:         cfg,
    gamechan:     gamechan,
    SensorChan:   sensorchan,
    led:          NewSensorLed(cfg),
    hit:          NewSensorHitInput(cfg, sensorchan),
    Logger:       log.New(logger.Writer(), fmt.Sprintf("[sensor:%s]: ", cfg.Id), logger.Flags()),
  }
}

func (s *Sensor) Start(parentctx context.Context) error {
  s.Printf("starting sensor")
  g, ctx := errgroup.WithContext(parentctx)

  if !s.IsTestSensor() {
    if err := s.led.Connect(); err != nil {
      s.Printf("error connecting to sensor: %s", err)
      return err
    }
  }

  g.Go(func() error {
    for {
      select {
      case evt := <-s.SensorChan:
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
  })

  return g.Wait()
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
