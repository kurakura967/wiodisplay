//go:build !tinygo

// Package hardware provides a koebiten-compatible Hardware implementation
// for Wio Terminal emulation via wiodisplay.
//
// On desktop (!tinygo), this package exposes WioHardware and WioDisplay types
// that satisfy the same method signatures as koebiten.Hardware and koebiten.Displayer.
// They are intended to be used only in tinygo builds; on desktop the types are
// provided as stubs so that `go build ./...` succeeds.
//
// Key mapping from Wio Terminal buttons to koebiten.Key indices:
//
//	Key0  (0) ← WIO_KEY_A   (Z key on desktop)
//	Key1  (1) ← WIO_KEY_B   (X key on desktop)
//	Key2  (2) ← WIO_KEY_C   (C key on desktop)
//	Key3  (3) ← WIO_5S_PRESS (Enter key on desktop)
//	KeyLeft  (16) ← WIO_5S_LEFT  (← arrow on desktop)
//	KeyRight (17) ← WIO_5S_RIGHT (→ arrow on desktop)
//	KeyUp    (18) ← WIO_5S_UP    (↑ arrow on desktop)
//	KeyDown  (19) ← WIO_5S_DOWN  (↓ arrow on desktop)
package hardware

import (
	"image/color"

	"github.com/kurakura967/wiodisplay/client"
	"github.com/kurakura967/wiodisplay/initdisplay"
)

// Device は wiodisplay のデスクトップ向け Hardware 実装のシングルトン。
// koebiten.SetHardware(hardware.Device) として渡すことを想定しているが、
// デスクトップでは koebiten 自体が //go:build tinygo のため、
// 実際の SetHardware 呼び出しは tinygo ビルド時のみ有効になる。
var Device = &WioHardware{}

// WioHardware はデスクトップ向けの Hardware スタブ実装。
type WioHardware struct {
	display *WioDisplay
}

// Init は RPC サーバー（wio-emu）への接続を確立し、表示デバイスを初期化する。
func (h *WioHardware) Init() error {
	d := initdisplay.InitDisplay()
	h.display = &WioDisplay{d: d}
	return nil
}

// GetDisplay は koebiten.Displayer 互換の WioDisplay を返す。
// デスクトップでは koebiten.Displayer 型は存在しないため、
// *WioDisplay を直接返す。tinygo ビルドとシグネチャが異なるが、
// このファイルは tinygo ビルドでは除外される。
func (h *WioHardware) GetDisplay() *WioDisplay {
	return h.display
}

// KeyUpdate はデスクトップスタブのため何もしない。
// デスクトップでは koebiten の RunGame は使わず、
// サーバー（server.go）が Ebiten のキー入力を直接処理する。
func (h *WioHardware) KeyUpdate() error {
	return nil
}

// WioDisplay はデスクトップ向けの Displayer スタブ実装。
// RPC 経由でサーバー（wio-emu）に描画コマンドを送信する。
type WioDisplay struct {
	d displayer
}

// displayer は WioDisplay が必要とする描画メソッドの内部インターフェース。
type displayer interface {
	Size() (x, y int16)
	SetPixel(x, y int16, c color.RGBA)
	Display() error
	FillScreen(c color.RGBA)
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

// ClearBuffer は画面全体を黒でクリアする。
func (w *WioDisplay) ClearBuffer() {
	w.d.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 255})
}

// ClearDisplay は画面表示をクリアする。
func (w *WioDisplay) ClearDisplay() {
	w.d.FillScreen(color.RGBA{R: 0, G: 0, B: 0, A: 255})
}

// GetButtonState はデスクトップで RPC 経由でボタン状態を取得する。
// server.Buttons の状態を読み取る。
func GetButtonState(pin int) bool {
	if client.Conn == nil {
		return false
	}
	var pressed bool
	client.Conn.Call("DisplayService.GetButtonState", &client.GetButtonStateArgs{Pin: pin}, &pressed)
	return pressed
}
