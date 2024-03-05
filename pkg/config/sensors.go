package config

import (
  "fmt"
  "log"
  "strings"
  "gopkg.in/yaml.v2"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type SensorConfig struct {
  Id            string
  Device        string
  Gpiochip      string
  Hitpin        string
  Ledpin        string
  Redpin        string
  Greenpin      string
  Bluepin       string
  Debounce      int
}

type SensorsConfig struct {
  Configs       map[string]*SensorConfig    `yaml:"configs" json:"configs"`
}

func NewSensorConfig(id, device, chip, hitpin, ledrgbpin string, debouncetime int) *SensorConfig {
  ledpin := ""
  redpin := ""
  greenpin := ""
  bluepin := ""

  parts := strings.Split(ledrgbpin, constants.SPLIT)
  if len(parts) == 3 { // rgb led
    redpin = parts[0]
    greenpin = parts[1]
    bluepin = parts[2]
  } else {
    ledpin = ledrgbpin // assumes single led pin
  }

  return &SensorConfig{
    Id:           id,
    Device:       device,
    Gpiochip:     chip,
    Hitpin:       hitpin,
    Ledpin:       ledpin,
    Redpin:       redpin,
    Greenpin:     greenpin,
    Bluepin:      bluepin,
    Debounce:     debouncetime,
  }
}

func DefaultSensorConfig(id string) *SensorConfig {
  return &SensorConfig{
    Id:           id,
    Device:       "",
    Gpiochip:     "gpiochip0",
    Hitpin:       "",
    Ledpin:       "",
    Redpin:       "",
    Greenpin:     "",
    Bluepin:      "",
    Debounce:     100,
  }
}

func NewSensorsConfig() *SensorsConfig {
  return &SensorsConfig{
    Configs:    map[string]*SensorConfig{},
  }
}

func (sc *SensorConfig) IsRGB() bool {
  return sc.Redpin != "" && sc.Greenpin != "" && sc.Bluepin != ""
}

func (sc *SensorConfig) Enabled() bool {
  return sc.Device != ""
}

func (sc *SensorConfig) Error() error {
  if strings.HasPrefix(sc.Id, constants.TEST_SENSOR_PREFIX) {
    return constants.ERR_TEST_SENSOR
  }
  if sc.Device == "" {
    return constants.ERR_NO_SENSOR_DEVICE
  }
  if sc.Gpiochip == "" {
    return constants.ERR_NO_SENSOR_GPIOCHIP
  }
  if sc.Hitpin == "" {
    return constants.ERR_NO_SENSOR_HITPIN
  }
  if sc.Ledpin == "" {
    return constants.ERR_NO_SENSOR_LEDPIN
  }
  return nil
}

func (sc *SensorsConfig) Set(value string) error {
  parts := strings.SplitN(value, constants.SPLIT, 5) // split at most into 5 parts
  id := parts[0]
  log.Printf("parsing sensor %s", value)

  if strings.HasPrefix(id, constants.TEST_SENSOR_PREFIX) {
    sc.Configs[id] = DefaultSensorConfig(id)
    return nil
  }

  if len(parts) != 5 {
    return constants.ERR_INVALID_SENSOR_FLAG
  }

  dev := parts[1]
  chip := parts[2]
  hit := parts[3]
  led := parts[4] // either a single pin or 3 rgb pins, i.e. 'gpio13' or 'gpio13:gpio14:gpio17' for rgb led

  sc.Configs[id] = NewSensorConfig(id, dev, chip, hit, led, 100)
  return nil
}

func (sc *SensorsConfig) String() string {
  yamlBytes, err := yaml.Marshal(sc)
  if err != nil {
    return fmt.Sprintf("Error marshalling sensors config into yaml: %s", err)
  }
  return string(yamlBytes)
}
