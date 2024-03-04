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
  Debounce      int
}

type SensorsConfig struct {
  Configs       map[string]*SensorConfig    `yaml:"configs" json:"configs"`
}

func NewSensorConfig(id, device, chip, hitpin, ledpin string, debouncetime int) *SensorConfig {
  return &SensorConfig{
    Id:           id,
    Device:       device,
    Gpiochip:     chip,
    Hitpin:       hitpin,
    Ledpin:       ledpin,
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
    Debounce:     100,
  }
}

func NewSensorsConfig() *SensorsConfig {
  return &SensorsConfig{
    Configs:    map[string]*SensorConfig{},
  }
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
  parts := strings.Split(value, constants.SPLIT)
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
  led := parts[4]

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
