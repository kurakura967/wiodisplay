//go:build tinygo

package initdisplay

import (
	"github.com/kurakura967/wiodisplay/driver/ili9341"
	"github.com/kurakura967/wiodisplay/machine"
)

// InitDisplay は実機の ili9341 ドライバを初期化して返す。
func InitDisplay() *ili9341.Device {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40000000,
	})

	backlight := machine.LCD_BACKLIGHT
	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})

	display := ili9341.NewSPI(machine.SPI3, machine.LCD_DC, machine.LCD_SS_PIN, machine.LCD_RESET)
	display.Configure(ili9341.Config{})
	backlight.High()
	display.SetRotation(ili9341.Rotation270)

	return display
}
