package sensor

import (
  "log"
  "sync"
  "time"
  "context"
  "golang.org/x/sync/errgroup"
  "github.com/taemon1337/gpiod"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/game"
)

type SensorHitInput struct {
  conf          *config.SensorConfig    `yaml:"config" json:"config"`
  color         string                  `yaml:"color" json:"color"`
  line          *gpiod.Line             `yaml:"-" json:"-"`
  hitchan       chan gpiod.LineEvent    `yaml:"-" json:"-"`
  HitChan       chan game.GameEvent     `yaml:"-" json:"-"`
  lasthit       time.Time               `yaml:"last_hit" json:"last_hit"`
  lock          *sync.Mutex             `yaml:"-" json:"-"`
  *log.Logger
}

func NewSensorHitInput(cfg *config.SensorConfig, logger *log.Logger) *SensorHitInput {
  return &SensorHitInput{
    conf:           cfg,
    color:          "",
    line:           nil,
    hitchan:        make(chan gpiod.LineEvent, constants.CHANNEL_WIDTH),
    HitChan:        make(chan game.GameEvent, constants.CHANNEL_WIDTH),
    lasthit:        time.Now(),
    lock:           &sync.Mutex{},
    Logger:   logger,
  }
}

func (s *SensorHitInput) ProcessEvent(evt gpiod.LineEvent) {
  debounce_duration := time.Duration(s.conf.Debounce) * time.Millisecond

  if time.Since(s.lasthit) < debounce_duration {
    s.Printf("IGNORING DUP HIT: %s", evt)
    return // ignore since within debounce window
  }

  s.Printf("HIT: %s", evt)

  s.lock.Lock()
  s.lasthit = time.Now() // hittime is last debounced hit time
  s.lock.Unlock()

  select {
    case s.HitChan <- game.NewGameEvent(constants.SENSOR_HIT, []byte("1")):
    default:
      s.Printf("hit channel is full - discarding event: %s", evt)
  }
}

func (s *SensorHitInput) Start(parentctx context.Context) error {
  hitpin, err := ParseGpioPin(s.conf.Device, s.conf.Hitpin)
  if err != nil {
    return err
  }

  s.Printf("gpio hit pin: %d", hitpin)

  // event channel buffer
  eh := func(evt gpiod.LineEvent) {
    select {
    case s.hitchan <- evt:
    default:
      s.Printf("event chan overflow - discarding event")
    }
  }

  hit, err := gpiod.RequestLine(s.conf.Gpiochip, hitpin, gpiod.WithPullUp, gpiod.WithRisingEdge, gpiod.WithEventHandler(eh))
  if err != nil {
    s.Printf("cannot request gpiod %d hit line: %s", hitpin, err)
    return err
  }

  s.line = hit

  g, ctx := errgroup.WithContext(parentctx)

  g.Go(func() error {
    defer func() {
      hit.Reconfigure(gpiod.AsInput)
      hit.Close()
    }()

    for {
      select {
      case evt := <-s.hitchan:
        s.ProcessEvent(evt)
      case <-ctx.Done():
        s.Printf("stopping hit input sensor")
        return ctx.Err()
      }
    }
  })

  return g.Wait()
}
