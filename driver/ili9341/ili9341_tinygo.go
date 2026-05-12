//go:build tinygo

package ili9341

import (
	real "tinygo.org/x/drivers/ili9341"

	"github.com/kurakura967/wiodisplay/machine"
)

// 実機では tinygo.org/x/drivers/ili9341 をそのまま使用する。

type Rotation = real.Rotation
type Config = real.Config
type Device = real.Device

const Rotation270 = real.Rotation270

func NewSPI(bus machine.SPI, dc, cs, rst machine.Pin) *Device {
	return real.NewSPI(bus, dc, cs, rst)
}
