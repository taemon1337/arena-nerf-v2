package sensor

import (
  "log"
  "sync"
  "time"
  "github.com/taemon1337/gpiod"
  "github.com/taemon1337/arena-nerf/pkg/constants"
  "github.com/taemon1337/arena-nerf/pkg/config"
)

type SensorLed struct {
  conf          *config.SensorConfig    `yaml:"config" json:"config"`
  color         string                  `yaml:"color" json:"color"`
  line          *gpiod.Line             `yaml:"-" json:"-"`
  lock          *sync.Mutex             `yaml:"-" json:"-"`
  *log.Logger
}

func NewSensorLed(cfg *config.SensorConfig, logger *log.Logger) *SensorLed {
  return &SensorLed{
    conf:     cfg,
    color:    "",
    line:     nil,
    lock:     &sync.Mutex{},
    Logger:   logger,
  }
}

func (led *SensorLed) Connect() error {
  ledpin, err := ParseGpioPin(led.conf.Device, led.conf.Ledpin)
  if err != nil {
    led.Printf("cannot parse gpio led pin %s: %s", ledpin, err)
    return err
  }

  ledline, err := gpiod.RequestLine(led.conf.Gpiochip, ledpin, gpiod.AsOutput(constants.OFF))
  if err != nil {
    led.Printf("cannot request gpiod %d led line: %s", ledpin, err)
    return err
  }

  led.line = ledline

  // start by blinking led
  log.Printf("Blinking LED 5 times...")
  time.Sleep(3 * time.Second)
  led.Blink(5)
  time.Sleep(1 * time.Second)
  return nil
}

func (led *SensorLed) Close() {
  led.line.Reconfigure(gpiod.AsInput)
  led.line.Close()
}

func (led *SensorLed) SetColor(color string) {
  led.lock.Lock()
  defer led.lock.Unlock()
  led.color = color
}

func (led *SensorLed) GetColor() string {
  return led.color
}

func (led *SensorLed) Blink(times int) {
  for i := 0; i < times; i++ {
    led.BlinkOnce()
    time.Sleep(constants.BLINK_DELAY)
  }
}

func (led *SensorLed) BlinkOnce() {
  if led.line != nil {
    led.line.SetValue(constants.ON)
    time.Sleep(constants.BLINK_DELAY)
    led.line.SetValue(constants.OFF)
  }
}
