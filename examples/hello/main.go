package main

import (
	"image/color"
	"time"

	"github.com/kurakura967/wiodisplay/initdisplay"
	"github.com/kurakura967/wiodisplay/machine"
)

func main() {
	display := initdisplay.InitDisplay()

	// 画面を青でクリア
	display.FillScreen(color.RGBA{R: 0, G: 0, B: 128, A: 255})

	// 中央に白い矩形
	display.FillRectangle(60, 80, 200, 80, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	for {
		// ボタン A(Z キー) を押すと赤、B(X キー) を押すと緑、C(C キー) を押すと青に塗り替え
		if !machine.WIO_KEY_A.Get() {
			display.FillScreen(color.RGBA{R: 200, G: 0, B: 0, A: 255})
		} else if !machine.WIO_KEY_B.Get() {
			display.FillScreen(color.RGBA{R: 0, G: 200, B: 0, A: 255})
		} else if !machine.WIO_KEY_C.Get() {
			display.FillScreen(color.RGBA{R: 0, G: 0, B: 200, A: 255})
		}

		time.Sleep(16 * time.Millisecond)
	}
}
