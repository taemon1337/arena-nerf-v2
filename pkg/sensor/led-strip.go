package sensor

import (
  "log"
  "time"
  "github.com/rpi-ws281x/rpi-ws281x-go"
  "github.com/taemon1337/arena-nerf/pkg/config"
  "github.com/taemon1337/arena-nerf/pkg/constants"
)

type RGB struct {
	R, G, B uint8
}

func (c RGB) Value() uint32 {
  return uint32(c.R)<<16 | uint32(c.G)<<8 | uint32(c.B)
}

// LedStrip represents an LED strip with RGB colors and GPIO lines for communication
type LedStrip struct {
  conf          *config.SensorConfig
  numLEDs       int
  ledstrip      *ws2811.WS2811
  datapin       string
  *log.Logger
}

// NewLedStrip initializes a new LedStrip instance
func NewLedStrip(cfg *config.SensorConfig, logger *log.Logger) *LedStrip {
  return &LedStrip{
    conf:       cfg,
    numLEDs:    cfg.Ledcount,
    ledstrip:   nil,
    datapin:    cfg.Ledpin,
    Logger:     logger,
  }
}

func (strip *LedStrip) Connect() error {
  datapin, err := ParseGpioPin(strip.conf.Device, strip.datapin)
  if err != nil {
    log.Printf("cannot parse led strip data pin %s: %s", strip.datapin, err)
    return err
  }

  strip.Printf("gpio LED pin: %d", datapin)

  width := 1
  height := strip.numLEDs
  bright := 64 // 0-255
  size := width * height
  opt := ws2811.DefaultOptions
  opt.Channels[0].Brightness = bright
  opt.Channels[0].LedCount = size
  opt.Channels[0].GpioPin = datapin

  ws, err := ws2811.MakeWS2811(&opt)
  if err != nil {
    strip.Printf("could not get new LED device: %s", err)
    return err
  }

  err = ws.Init()
  if err != nil {
    strip.Printf("could not initialize LED strip: %s", err)
    return err
  }

  strip.ledstrip = ws

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
  return strip.ledstrip != nil
}

func (strip *LedStrip) Close() {
  strip.ledstrip.Fini()
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
  for i := 0; i < len(strip.ledstrip.Leds(0)); i++ {
    strip.ledstrip.Leds(0)[i] = color.Value()
  }
}

// On turns all LEDs to the same given color
func (strip *LedStrip) On(color RGB) error {
  for i := 0; i < strip.numLEDs; i++ {
    strip.SetLEDColor(i, color)
  }
  return strip.ledstrip.Render()
}

// Off turns off all LEDs by setting their colors to black
func (strip *LedStrip) Off() error {
  black := RGB{0, 0, 0}
  return strip.On(black)
}
