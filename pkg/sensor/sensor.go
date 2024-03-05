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
  enableLeds    bool
  enableHits    bool
  *log.Logger
}

func NewSensor(cfg *config.SensorConfig, gamechan *game.GameChannel, logger *log.Logger, enable_leds, enable_hits bool) *Sensor {
  sensorchan := make(chan game.GameEvent, constants.CHANNEL_WIDTH)
  logger = log.New(logger.Writer(), fmt.Sprintf("[sensor:%s]: ", cfg.Id), logger.Flags())

  return &Sensor{
    conf:         cfg,
    gamechan:     gamechan,
    SensorChan:   sensorchan,
    led:          NewSensorLed(cfg, logger),
    hit:          NewSensorHitInput(cfg, sensorchan, logger),
    enableLeds:   enable_leds,
    enableHits:   enable_hits,
    Logger:       logger,
  }
}

func (s *Sensor) Start(parentctx context.Context) error {
  s.Printf("starting sensor %s", s.conf.Id)
  g, ctx := errgroup.WithContext(parentctx)

  if s.LedEnabled() {
    if err := s.led.Connect(); err != nil {
      s.Printf("error connecting to sensor: %s", err)
      return err
    }

    defer s.led.Close()
  }

  if s.HitEnabled() {
    g.Go(func() error {
      return s.hit.Start(ctx)
    })
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
        s.Printf("stopping sensor %s", s.conf.Id)
        return ctx.Err()
      default:
        time.Sleep(3 * time.Second) // do something later
      }
    }
  })

  return g.Wait()
}

func (s *Sensor) SensorHit(sensorid string) {
  if !s.IsTestSensor() {
    s.led.Blink(5)
  }

  pay := strings.Join([]string{sensorid, s.led.GetColor(), "1"}, constants.SPLIT)
  s.gamechan.GameChan <- game.NewGameEvent(constants.SENSOR_HIT, []byte(pay))
}

func (s *Sensor) IsTestSensor() bool {
  return strings.HasPrefix(s.conf.Id, constants.TEST_SENSOR_PREFIX)
}

func (s *Sensor) LedEnabled() bool {
  return s.conf.Ledpin != "" && s.enableLeds && !s.IsTestSensor()
}

func (s *Sensor) HitEnabled() bool {
  return s.conf.Hitpin != "" && s.enableHits && !s.IsTestSensor()
}

func (s *Sensor) Led() *SensorLed {
  return s.led
}

func (s *Sensor) Hit() *SensorHitInput {
  return s.hit
}
