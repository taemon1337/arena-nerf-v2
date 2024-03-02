package constants

import (
  "errors"
)

var (
  TEST_SENSOR_PREFIX = "test"
  RANDOM_SENSOR_ID = "rand"
  RANDOM_COLOR_ID = "rand"
  NONE_SENSOR_ID = "none"
  ERR_SENSORS_DISABLED = errors.New("sensors are disabled")
  ERR_NO_SENSORS = errors.New("no sensors setup")
  ERR_NO_SENSOR_BY_NAME = errors.New("no sensor found by name")
  ERR_NO_SENSOR_DEVICE = errors.New("no sensor device set")
  ERR_NO_SENSOR_GPIOCHIP = errors.New("no sensor gpiochip set")
  ERR_NO_SENSOR_HITPIN = errors.New("no sensor hitpin set")
  ERR_NO_SENSOR_LEDPIN = errors.New("no sensor ledpin set")
  ERR_INVALID_SENSOR_FLAG = errors.New("invalid -sensor flag; expects <1-4>:<device>:<gpiochip>:<hit-pin>:<led-pin>")
  ERR_INVALID_SENSOR_NUMBER = errors.New("invalid -sensor <number>, must be 1-4")
  ERR_TEST_SENSOR = errors.New("sensor is a test only sensor")
)

