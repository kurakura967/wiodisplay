//go:build !tinygo

package ili9341

import (
	"image/color"

	"github.com/kurakura967/wiodisplay/client"
	"github.com/kurakura967/wiodisplay/machine"
)

// Rotation represents display rotation.
type Rotation uint8

const Rotation270 Rotation = 3

// Config holds display configuration (unused in emulator).
type Config struct{}

// Device is the ili9341 display that forwards drawing calls to the server via RPC.
type Device struct{}

func colorToRGB565(c color.RGBA) uint16 {
	r := uint16(c.R >> 3)
	g := uint16(c.G >> 2)
	b := uint16(c.B >> 3)
	return (r << 11) | (g << 5) | b
}

// NewSPI creates a new Device. Arguments are accepted for API compatibility with
// the real driver but are ignored on desktop.
func NewSPI(bus machine.SPI, dc, cs, rst machine.Pin) *Device {
	return &Device{}
}

func (d *Device) Configure(c Config) {}

func (d *Device) SetRotation(r Rotation) {}

// SetPixel satisfies the tinyfont.Displayer interface.
func (d *Device) SetPixel(x, y int16, c color.RGBA) {
	d.DrawPixel(x, y, c)
}

func (d *Device) DrawPixel(x, y int16, c color.RGBA) {
	var reply struct{}
	client.Conn.Call("DisplayService.DrawPixel", &client.DrawPixelArgs{
		X: x, Y: y, Color: colorToRGB565(c),
	}, &reply)
}

func (d *Device) FillRectangle(x, y, w, h int16, c color.RGBA) {
	var reply struct{}
	client.Conn.Call("DisplayService.FillRectangle", &client.FillRectangleArgs{
		X: x, Y: y, W: w, H: h, Color: colorToRGB565(c),
	}, &reply)
}

func (d *Device) FillScreen(c color.RGBA) {
	var reply struct{}
	client.Conn.Call("DisplayService.FillScreen", &client.FillScreenArgs{
		Color: colorToRGB565(c),
	}, &reply)
}

// Display satisfies the drivers.Displayer interface (no-op for emulator).
func (d *Device) Display() error { return nil }

// Size returns the display dimensions.
func (d *Device) Size() (x, y int16) {
	return 320, 240
}
