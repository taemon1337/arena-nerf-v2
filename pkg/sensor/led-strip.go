package sensor

import (
  "log"
  "time"
  "github.com/taemon1337/gpiod"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type RGB struct {
	R, G, B uint8
}

// LedStrip represents an LED strip with RGB colors and GPIO lines for communication
type LedStrip struct {
  conf          *config.SensorConfig
  numLEDs       int
  ledStrip      []RGB
  datapin       string
  dataline      *gpiod.Line
  *log.Logger
}

// NewLedStrip initializes a new LedStrip instance
func NewLedStrip(cfg *config.SensorConfig, logger *log.Logger) *LedStrip {
  numleds := cfg.Ledcount
  ledstrip := make([]RGB, numleds)
  for i := range ledstrip {
    ledstrip[i] = RGB{0, 0, 0} // black|off
  }

  return &LedStrip{
    conf:       cfg,
    numLEDs:    numleds,
    ledStrip:   ledstrip,
    datapin:    cfg.Ledpin,
    dataline:   nil,
    Logger:     logger,
  }
}

func (strip *LedStrip) Connect() error {
  datapin, err := ParseGpioPin(strip.conf.Device, strip.datapin)
  if err != nil {
    log.Printf("cannot parse led strip data pin %s: %s", strip.datapin, err)
    return err
  }

  dataline, err := gpiod.RequestLine(strip.conf.Gpiochip, datapin, gpiod.AsOutput(constants.OFF))
  if err != nil {
    log.Printf("cannot request gpiod %d led data line: %s", datapin, err)
    return err
  }

  strip.dataline = dataline

  // start by blinking led
  log.Printf("Blinking LED strip 5 times...")
  time.Sleep(3 * time.Second)
  if err := strip.Blink(5, RGB{255, 255, 255}); err != nil {
    return err
  }
  time.Sleep(1 * time.Second)
  return nil
}

func (strip *LedStrip) Connected() bool {
  return strip.dataline != nil
}

// On turns all LEDs to the same given color
func (strip *LedStrip) Blink(times int, color RGB) error {
  for i := 0; i < times; i++ {
    if err := strip.BlinkOnce(color); err != nil {
      return err
    }
    time.Sleep(constants.BLINK_DELAY)
  }
  return nil
}

func (strip *LedStrip) BlinkOnce(color RGB) error {
  if strip.Connected() {
    if err := strip.On(color); err != nil {
      return err
    }
    time.Sleep(constants.BLINK_DELAY)
    if err := strip.Off(); err != nil {
      return err
    }
  } else {
    strip.Printf("strip not connected")
  }
  return nil
}

// SetLEDColor sets the color of an individual LED by its index
func (strip *LedStrip) SetLEDColor(index int, color RGB) {
  if index >= 0 && index < strip.numLEDs {
    strip.ledStrip[index] = color
  }
}

// On turns all LEDs to the same given color
func (strip *LedStrip) On(color RGB) error {
  for i := range strip.ledStrip {
    strip.SetLEDColor(i, color)
  }
  return strip.SendData()
}

// Off turns off all LEDs by setting their colors to black
func (strip *LedStrip) Off() error {
  black := RGB{0, 0, 0}
  return strip.On(black)
}

// SendData sends data to NeoPixels
func (strip *LedStrip) SendData() error {
  strip.Printf("rendering leds")
  // Send data to NeoPixels bit by bit
  for _, rgb := range strip.ledStrip {
      // Send RGB data
      strip.sendRGBData(rgb.R, rgb.G, rgb.B)
  }

  return nil
}

// SendRGBData sends RGB data to NeoPixels
func (strip *LedStrip) sendRGBData(red, green, blue uint8) {
  for i := 7; i >= 0; i-- {
    strip.sendDataBit((green >> uint(i)) & 0x01)
  }
  for i := 7; i >= 0; i-- {
    strip.sendDataBit((red >> uint(i)) & 0x01)
  }
  for i := 7; i >= 0; i-- {
    strip.sendDataBit((blue >> uint(i)) & 0x01)
  }
}

// SendDataBit sends a single bit of data to NeoPixels
func (strip *LedStrip) sendDataBit(bit uint8) {
  if bit == 1 {
    if err := strip.dataline.SetValue(1); err != nil {
      // Handle error
      return
    }
  } else {
    if err := strip.dataline.SetValue(0); err != nil {
      // Handle error
      return
    }
  }
}

