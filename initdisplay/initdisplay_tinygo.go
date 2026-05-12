//go:build tinygo

package initdisplay

import (
	"github.com/kurakura967/wiodisplay/driver/ili9341"
	"github.com/kurakura967/wiodisplay/machine"
)

// InitDisplay は実機の ili9341 ドライバを初期化して返す。
func InitDisplay() *ili9341.Device {
	return ili9341.NewSPI(machine.SPI3, machine.LCD_DC, machine.LCD_SS_PIN, machine.LCD_RESET)
}
