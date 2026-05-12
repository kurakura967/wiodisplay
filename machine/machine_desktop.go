//go:build !tinygo

package machine

import "github.com/kurakura967/wiodisplay/client"

// SPI はスタブ。API 互換のため定義する。
type SPI struct{}

type SPIConfig struct {
	SCK       Pin
	SDO       Pin
	SDI       Pin
	Frequency uint32
}

func (s SPI) Configure(c SPIConfig) {}
func (s SPI) Tx(w, r []byte) error  { return nil }

var SPI3 SPI

// Pin はデスクトップでは RPC 経由でサーバーに問い合わせる GPIO ピンスタブ。
type Pin uint8

type PinMode uint8

const (
	PinInput       PinMode = 0
	PinInputPullup PinMode = 1
	PinOutput      PinMode = 2
)

type PinConfig struct {
	Mode PinMode
}

func (p Pin) Configure(c PinConfig) {}
func (p Pin) High()                 {}
func (p Pin) Low()                  {}

// Get は TCP 経由でサーバーのボタン状態を取得する。
// Wio Terminal のボタンはアクティブローなので、押下時は false を返す。
func (p Pin) Get() bool {
	if client.Conn == nil {
		return true // 未接続時は非押下
	}
	var pressed bool
	client.Conn.Call("DisplayService.GetButtonState", &client.GetButtonStateArgs{Pin: int(p)}, &pressed)
	return !pressed
}

// LCD ピン定数
const (
	LCD_SCK_PIN   Pin = 0
	LCD_SDO_PIN   Pin = 1
	LCD_SDI_PIN   Pin = 2
	LCD_DC        Pin = 3
	LCD_SS_PIN    Pin = 4
	LCD_RESET     Pin = 5
	LCD_BACKLIGHT Pin = 6
)

// Wio Terminal ボタン・ジョイスティックピン定数
const (
	WIO_KEY_A    Pin = 10
	WIO_KEY_B    Pin = 11
	WIO_KEY_C    Pin = 12
	WIO_5S_UP    Pin = 13
	WIO_5S_DOWN  Pin = 14
	WIO_5S_LEFT  Pin = 15
	WIO_5S_RIGHT Pin = 16
	WIO_5S_PRESS Pin = 17
)
