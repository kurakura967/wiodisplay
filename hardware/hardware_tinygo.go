//go:build tinygo

package hardware

import (
	"image/color"
	"machine"

	"github.com/sago35/koebiten"
	"github.com/kurakura967/wiodisplay/driver/ili9341"
)

// WioHardware は Wio Terminal 実機向けの koebiten.Hardware 実装。
// tinygo build 時のみ有効。
var Device = &WioHardware{}

type WioHardware struct {
	display *WioDisplay
	pins    []machine.Pin
	state   []pinState
	keybuf  [1]koebiten.Key
}

type pinState uint8

const (
	stateNone pinState = iota
	stateNoneToPress
	statePress
	statePressToRelease
)

// Init は SPI・バックライト・ili9341 ディスプレイを初期化する。
func (h *WioHardware) Init() error {
	machine.SPI3.Configure(machine.SPIConfig{
		SCK:       machine.LCD_SCK_PIN,
		SDO:       machine.LCD_SDO_PIN,
		SDI:       machine.LCD_SDI_PIN,
		Frequency: 40000000,
	})

	backlight := machine.LCD_BACKLIGHT
	backlight.Configure(machine.PinConfig{Mode: machine.PinOutput})

	d := ili9341.NewSPI(
		machine.SPI3,
		machine.LCD_DC,
		machine.LCD_SS_PIN,
		machine.LCD_RESET,
	)
	d.Configure(ili9341.Config{})
	d.SetRotation(ili9341.Rotation270)
	d.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 255})

	backlight.High()

	h.display = &WioDisplay{d: d}

	// gpioPins のインデックスが koebiten.Key の値に対応する。
	// NoPin のインデックスはそのキーが存在しないことを意味する。
	configPins := []machine.Pin{
		machine.WIO_KEY_A,
		machine.WIO_KEY_B,
		machine.WIO_KEY_C,
		machine.WIO_5S_PRESS,
		machine.WIO_5S_UP,
		machine.WIO_5S_LEFT,
		machine.WIO_5S_DOWN,
		machine.WIO_5S_RIGHT,
	}
	for _, p := range configPins {
		p.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	}

	// koebiten.Key インデックスに合わせて並べる:
	//   Key0=WIO_KEY_A, Key1=WIO_KEY_B, Key2=WIO_KEY_C, Key3=WIO_5S_PRESS,
	//   [4..15]=NoPin,
	//   KeyLeft(16)=WIO_5S_LEFT, KeyRight(17)=WIO_5S_RIGHT,
	//   KeyUp(18)=WIO_5S_UP, KeyDown(19)=WIO_5S_DOWN
	h.pins = []machine.Pin{
		machine.WIO_KEY_A,    // 0 = Key0
		machine.WIO_KEY_B,    // 1 = Key1
		machine.WIO_KEY_C,    // 2 = Key2
		machine.WIO_5S_PRESS, // 3 = Key3 (ジョイスティック中央)
		machine.NoPin,        // 4
		machine.NoPin,        // 5
		machine.NoPin,        // 6
		machine.NoPin,        // 7
		machine.NoPin,        // 8
		machine.NoPin,        // 9
		machine.NoPin,        // 10
		machine.NoPin,        // 11
		machine.NoPin,        // 12
		machine.NoPin,        // 13
		machine.NoPin,        // 14
		machine.NoPin,        // 15
		machine.WIO_5S_LEFT,  // 16 = KeyLeft
		machine.WIO_5S_RIGHT, // 17 = KeyRight
		machine.WIO_5S_UP,    // 18 = KeyUp
		machine.WIO_5S_DOWN,  // 19 = KeyDown
	}
	h.state = make([]pinState, len(h.pins))
	return nil
}

// GetDisplay は koebiten.Displayer を実装した WioDisplay を返す。
func (h *WioHardware) GetDisplay() koebiten.Displayer {
	return h.display
}

// KeyUpdate はボタン状態をポーリングして koebiten の入力状態を更新する。
func (h *WioHardware) KeyUpdate() error {
	buf := h.keybuf[:]
	for i, pin := range h.pins {
		pressed := false
		if pin != machine.NoPin {
			pressed = !pin.Get() // アクティブロー
		}

		switch h.state[i] {
		case stateNone:
			if pressed {
				h.state[i] = stateNoneToPress
			}
		case stateNoneToPress:
			h.state[i] = statePress
			buf[0] = koebiten.Key(i)
			koebiten.AppendJustPressedKeys(buf)
		case statePress:
			buf[0] = koebiten.Key(i)
			koebiten.AppendPressedKeys(buf)
			if !pressed {
				h.state[i] = statePressToRelease
			}
		case statePressToRelease:
			h.state[i] = stateNone
			buf[0] = koebiten.Key(i)
			koebiten.AppendJustReleasedKeys(buf)
		}
	}
	return nil
}

// WioDisplay は ili9341.Device をラップして koebiten.Displayer を実装する。
type WioDisplay struct {
	d *ili9341.Device
}

func (w *WioDisplay) Size() (x, y int16) {
	return w.d.Size()
}

func (w *WioDisplay) SetPixel(x, y int16, c color.RGBA) {
	w.d.SetPixel(x, y, c)
}

func (w *WioDisplay) Display() error {
	return w.d.Display()
}

// ClearBuffer は no-op。
// koebiten の RunGame は毎フレーム ClearBuffer を呼ぶが、
// 実機（ILI9341/SPI）では全画面クリアがボトルネックになるため何もしない。
// 画面のクリアは Game.Draw() 側で dirty rect として管理する。
func (w *WioDisplay) ClearBuffer() {}

// ClearDisplay は画面表示をクリアする。
// ili9341 はバッファレス（直接描画）なため ClearBuffer と同じ動作をする。
func (w *WioDisplay) ClearDisplay() {
	w.d.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 255})
}
