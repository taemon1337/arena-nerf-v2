package sensor

import (
  "errors"
  "strconv"
  "strings"

  "github.com/taemon1337/gpiod/device/orangepi"
  "github.com/taemon1337/gpiod/device/rpi"
)

var (
  ERR_EMPTY_PIN = errors.New("no gpio pin specified")
  ERR_UNSUPPORTED_DEVICE = errors.New("unsupported sensor device")
)

func ParseGpioPin(device, pinstr string) (int, error) {
  if pinstr == "" {
    return 0, ERR_EMPTY_PIN
  }

  switch device {
    case "opi", "orangepi", "orangepi3", "orangepi3zero":
      return ParseOrangePiGpioPin(pinstr)
    case "rpi", "raspberrypi", "raspberrypi4", "raspberrypi4zero":
      return ParseRaspberryPiGpioPin(pinstr)
    default:
      return 0, ERR_UNSUPPORTED_DEVICE
  }
}

func ParseOrangePiGpioPin(pinstr string) (int, error) {
  pinstr = strings.ToLower(pinstr)
  
  switch pinstr {
    case "gpio2":
      return orangepi.GPIO2, nil
    case "gpio3":
      return orangepi.GPIO3, nil
    case "gpio4":
      return orangepi.GPIO4, nil
    case "gpio5":
      return orangepi.GPIO5, nil
    case "gpio6":
      return orangepi.GPIO6, nil
    case "gpio7":
      return orangepi.GPIO7, nil
    case "gpio8":
      return orangepi.GPIO8, nil
    case "gpio9":
      return orangepi.GPIO9, nil
    case "gpio10":
      return orangepi.GPIO10, nil
    case "gpio11":
      return orangepi.GPIO11, nil
    case "gpio12":
      return orangepi.GPIO12, nil
    case "gpio13":
      return orangepi.GPIO13, nil
    case "gpio14":
      return orangepi.GPIO14, nil
    case "gpio15":
      return orangepi.GPIO15, nil
    case "gpio16":
      return orangepi.GPIO16, nil
    case "gpio17":
      return orangepi.GPIO17, nil
    default:
      return strconv.Atoi(pinstr)
  }
}

func ParseRaspberryPiGpioPin(pinstr string) (int, error) {
  pinstr = strings.ToLower(pinstr)
  
  switch pinstr {
    case "j8p27":
      return rpi.J8p27, nil
    case "j8p28":
      return rpi.J8p28, nil
    case "j8p3":
      return rpi.J8p3, nil
    case "j8p5":
      return rpi.J8p5, nil
    case "j8p7":
      return rpi.J8p7, nil
    case "j8p29":
      return rpi.J8p29, nil
    case "j8p31":
      return rpi.J8p31, nil
    default:
      return strconv.Atoi(pinstr)
  }
}
/*
  finish these pins later
	J8p26
	J8p24
	J8p21
	J8p19
	J8p23
	J8p32
	J8p33
	J8p8
	J8p10
	J8p36
	J8p11
	J8p12
	J8p35
	J8p38
	J8p40
	J8p15
	J8p16
	J8p18
	J8p22
	J8p37
	J8p13
*/

