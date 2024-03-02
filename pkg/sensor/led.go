package sensor

import (
  "sync"
)

type SensorLed struct {
  color         string          `yaml:"color" json:"color"`
  lock          *sync.Mutex     `yaml:"-" json:"-"`
}

func NewSensorLed(color string) *SensorLed {
  return &SensorLed{
    color:    color,
    lock:     &sync.Mutex{},
  }
}

func (led *SensorLed) SetColor(color string) {
  led.lock.Lock()
  defer led.lock.Unlock()
  led.color = color
}

func (led *SensorLed) GetColor() string {
  return led.color
}
