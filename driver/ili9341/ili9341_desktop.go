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

// DrawRGBBitmap8 draws a RGB565 big-endian image at (x, y).
func (d *Device) DrawRGBBitmap8(x, y int16, data []uint8, w, h int16) error {
	for dy := int16(0); dy < h; dy++ {
		for dx := int16(0); dx < w; dx++ {
			i := (dy*w + dx) * 2
			hi := uint16(data[i])
			lo := uint16(data[i+1])
			v := (hi << 8) | lo
			r := uint8((v>>11)&0x1F) << 3
			g := uint8((v>>5)&0x3F) << 2
			b := uint8(v&0x1F) << 3
			d.DrawPixel(x+dx, y+dy, color.RGBA{R: r, G: g, B: b, A: 0xFF})
		}
	}
	return nil
}

// Display satisfies the drivers.Displayer interface (no-op for emulator).
func (d *Device) Display() error { return nil }

// Size returns the display dimensions.
func (d *Device) Size() (x, y int16) {
	return 320, 240
}
