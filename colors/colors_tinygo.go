//go:build tinygo

// Package colors は Wio Terminal 向けの色生成ヘルパーを提供する。
//
// tinygo / desktop でそれぞれ異なる色型を返すが、呼び出し側は
// colors.RGB(r, g, b) という共通の書き方ができる。
package colors

import "tinygo.org/x/drivers/pixel"

// RGB は RGB 値から pixel.BaseColor（ILI9341 向け RGB565BE）を生成する。
func RGB(r, g, b uint8) pixel.BaseColor {
	return pixel.NewColor[pixel.RGB565BE](r, g, b)
}
