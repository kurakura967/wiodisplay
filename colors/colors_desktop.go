//go:build !tinygo

// Package colors は Wio Terminal 向けの色生成ヘルパーを提供する。
//
// tinygo / desktop でそれぞれ異なる色型を返すが、呼び出し側は
// colors.RGB(r, g, b) という共通の書き方ができる。
package colors

import "image/color"

// RGB は RGB 値から color.RGBA を生成する。A は常に 0xFF（不透明）。
func RGB(r, g, b uint8) color.RGBA {
	return color.RGBA{R: r, G: g, B: b, A: 0xFF}
}
