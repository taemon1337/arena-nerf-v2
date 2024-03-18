package sensor

import (
  "fmt"
  "log"
  "strings"
  "context"
  "golang.org/x/sync/errgroup"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

// https://github.com/hybridgroup/gobot/blob/release/drivers/gpio/rgb_led_driver.go

type Sensor struct {
  id            string
  conf          *config.SensorConfig
  gamechan      *game.GameChannel
  SensorChan    chan game.GameEvent
  led           *SensorLed
  ledstrip      *LedStrip
  hit           *SensorHitInput
  enableLeds    bool
  enableHits    bool
  *log.Logger
}

func NewSensor(id string, cfg *config.SensorConfig, gamechan *game.GameChannel, logger *log.Logger, enable_leds, enable_hits bool) *Sensor {
  logger.Printf("creating sensor %s", id)
  logger = log.New(logger.Writer(), fmt.Sprintf("[sensor:%s]: ", id), logger.Flags())

  return &Sensor{
    id:           id,
    conf:         cfg,
    gamechan:     gamechan,
    SensorChan:   make(chan game.GameEvent, constants.CHANNEL_WIDTH),
    led:          NewSensorLed(cfg, logger),
    ledstrip:     NewLedStrip(cfg, logger),
    hit:          NewSensorHitInput(cfg, logger),
    enableLeds:   enable_leds,
    enableHits:   enable_hits,
    Logger:       logger,
  }
}

func (s *Sensor) Start(parentctx context.Context) error {
  s.Printf("starting sensor %s", s.conf.Id)
  g, ctx := errgroup.WithContext(parentctx)

  if s.LedEnabled() {
    if s.LedStripEnabled() {
      if err := s.ledstrip.Connect(); err != nil {
        s.Printf("error connecting to LED strip on sensor %s: %s", s.conf.Id, err)
        return err
      }
    }

    if s.LedSingleEnabled() {
      if err := s.led.Connect(); err != nil {
        s.Printf("error connecting to LED on sensor %s: %s", s.conf.Id, err)
        return err
      }
    }
  }

  if s.HitEnabled() {
    g.Go(func() error {
      return s.hit.Start(ctx)
    })
  }

  g.Go(func() error {
    for {
      select {
      case evt := <-s.hit.HitChan:
        s.Printf("HIT CHAN: %s", evt)
        s.SensorHit(s.conf.Id)
        continue
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
        s.Close()
        return ctx.Err()
      }
    }
  })

  return g.Wait()
}

func (s *Sensor) Close() {
  if s.led.Connected() {
    s.led.Close()
  }
  if s.ledstrip.Connected() {
    s.ledstrip.Close()
  }
}

func (s *Sensor) SensorHit(sensorid string) {
  if !s.IsTestSensor() {
    s.led.Blink(5)
  }

  pay := strings.Join([]string{sensorid, s.led.GetColor(), "1"}, constants.SPLIT)
  select {
    case s.gamechan.GameChan <- game.NewGameEvent(constants.SENSOR_HIT, []byte(pay)):
    default:
      s.Printf("game chan is full - discarding event sensor hit")
  }
}

func (s *Sensor) IsTestSensor() bool {
  return strings.HasPrefix(s.conf.Id, constants.TEST_SENSOR_PREFIX)
}

func (s *Sensor) LedEnabled() bool {
  return s.conf.Ledpin != "" && s.enableLeds && !s.IsTestSensor()
}

func (s *Sensor) LedStripEnabled() bool {
  return s.conf.Ledpin != "" && s.conf.Ledcount > 1
}

func (s *Sensor) LedSingleEnabled() bool {
  return s.conf.Ledpin != "" && s.conf.Ledcount == 1
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
