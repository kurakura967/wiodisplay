//go:build tinygo

package machine

// 実機では TinyGo の組み込み machine パッケージをそのまま使用する。
// パッケージ名の衝突を避けるためエイリアスでインポートする。
import m "machine"

type Pin = m.Pin
type PinMode = m.PinMode
type PinConfig = m.PinConfig
type SPI = m.SPI
type SPIConfig = m.SPIConfig

const (
	PinInput       = m.PinInput
	PinInputPullup = m.PinInputPullup
	PinOutput      = m.PinOutput
)

var SPI3 = m.SPI3

const (
	LCD_SCK_PIN   = m.LCD_SCK_PIN
	LCD_SDO_PIN   = m.LCD_SDO_PIN
	LCD_SDI_PIN   = m.LCD_SDI_PIN
	LCD_DC        = m.LCD_DC
	LCD_SS_PIN    = m.LCD_SS_PIN
	LCD_RESET     = m.LCD_RESET
	LCD_BACKLIGHT = m.LCD_BACKLIGHT
)

const (
	WIO_KEY_A    = m.WIO_KEY_A
	WIO_KEY_B    = m.WIO_KEY_B
	WIO_KEY_C    = m.WIO_KEY_C
	WIO_5S_UP    = m.WIO_5S_UP
	WIO_5S_DOWN  = m.WIO_5S_DOWN
	WIO_5S_LEFT  = m.WIO_5S_LEFT
	WIO_5S_RIGHT = m.WIO_5S_RIGHT
	WIO_5S_PRESS = m.WIO_5S_PRESS
)
