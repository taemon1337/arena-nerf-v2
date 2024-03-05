package sensor

import (
  "log"
  "sync"
  "time"
  "github.com/taemon1337/gpiod"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type SensorRgbLed struct {
  conf          *config.SensorConfig    `yaml:"config" json:"config"`
  color         string                  `yaml:"color" json:"color"`
  redline       *gpiod.Line             `yaml:"-" json:"-"`
  greenline     *gpiod.Line             `yaml:"-" json:"-"`
  blueline      *gpiod.Line             `yaml:"-" json:"-"`
  lock          *sync.Mutex             `yaml:"-" json:"-"`
}

func NewSensorRgbLed(cfg *config.SensorConfig) *SensorRgbLed {
  return &SensorRgbLed{
    conf:       cfg,
    color:      "",
    redline:    nil,
    greenline:  nil,
    blueline:   nil,
    lock:       &sync.Mutex{},
  }
}

func (led *SensorRgbLed) Connect() error {
  redpin, err := ParseGpioPin(led.conf.Device, led.conf.Redpin)
  if err != nil {
    log.Printf("cannot parse gpio red led pin %s: %s", led.conf.Redpin, err)
    return err
  }

  greenpin, err := ParseGpioPin(led.conf.Device, led.conf.Greenpin)
  if err != nil {
    log.Printf("cannot parse gpio green led pin %s: %s", led.conf.Greenpin, err)
    return err
  }

  bluepin, err := ParseGpioPin(led.conf.Device, led.conf.Bluepin)
  if err != nil {
    log.Printf("cannot parse gpio blue led pin %s: %s", led.conf.Bluepin, err)
    return err
  }

  redline, err := gpiod.RequestLine(led.conf.Gpiochip, redpin, gpiod.AsOutput(constants.OFF))
  if err != nil {
    log.Printf("cannot request gpiod %d led line: %s", redpin, err)
    return err
  }

  greenline, err := gpiod.RequestLine(led.conf.Gpiochip, greenpin, gpiod.AsOutput(constants.OFF))
  if err != nil {
    log.Printf("cannot request gpiod %d led line: %s", greenpin, err)
    return err
  }

  blueline, err := gpiod.RequestLine(led.conf.Gpiochip, bluepin, gpiod.AsOutput(constants.OFF))
  if err != nil {
    log.Printf("cannot request gpiod %d led line: %s", bluepin, err)
    return err
  }

  led.redline = redline
  led.greenline = greenline
  led.blueline = blueline

  // start by blinking led
  log.Printf("Blinking LED 5 times...")
  time.Sleep(3 * time.Second)
  led.Blink(5)
  time.Sleep(1 * time.Second)
  return nil
}

func (led *SensorRgbLed) Close() {
  led.redline.Close()
  led.greenline.Close()
  led.blueline.Close()
}

func (led *SensorRgbLed) Blink(n int) {
  for i := 0; i < n; i++ {

  }
}


