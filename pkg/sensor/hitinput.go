package sensor

import (
  "log"
  "sync"
  "time"
  "context"
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
  SensorChan    chan game.GameEvent     `yaml:"-" json:"-"`
  lasthit       time.Time               `yaml:"last_hit" json:"last_hit"`
  lock          *sync.Mutex             `yaml:"-" json:"-"`
  *log.Logger
}

func NewSensorHitInput(cfg *config.SensorConfig, sensorchan chan game.GameEvent, logger *log.Logger) *SensorHitInput {
  return &SensorHitInput{
    conf:           cfg,
    color:          "",
    line:           nil,
    hitchan:        make(chan gpiod.LineEvent, constants.CHANNEL_WIDTH),
    SensorChan:     sensorchan,
    lasthit:        time.Now(),
    lock:           &sync.Mutex{},
    Logger:   logger,
  }
}

func (s *SensorHitInput) ProcessEvent(evt gpiod.LineEvent) {
  s.lock.Lock()
  defer s.lock.Unlock()

  debounce_duration := time.Duration(s.conf.Debounce) * time.Millisecond

  if time.Since(s.lasthit) < debounce_duration {
    return // ignore since within debounce window
  }

  s.lasthit = time.Now() // hittime is last debounced hit time
  s.SensorChan <- game.NewGameEvent(constants.SENSOR_HIT, []byte("1"))
}

func (s *SensorHitInput) Start(ctx context.Context) error {
  hitpin, err := ParseGpioPin(s.conf.Device, s.conf.Hitpin)
  if err != nil {
    return err
  }

  s.Printf("Sensor Hit input Hit pin: %d", hitpin)

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

  defer func() {
    hit.Reconfigure(gpiod.AsInput)
    hit.Close()
  }()

  done := false
  for !done {
    select {
    case evt := <-s.hitchan:
      s.Printf("HIT: %s", evt)
      s.ProcessEvent(evt)
    }
  }
  return constants.ERR_SENSOR_HIT_STOPPED
}

