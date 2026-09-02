//go:build !tinygo

// Package pins defines Pin 番号の定数。
// machine パッケージと server パッケージの両方から参照することで、
// Pin 値を1箇所で管理する。
package pins

// Wio Terminal のボタン・ジョイスティックに対応する Pin 番号。
// machine.Pin の値として machine パッケージで使われ、
// server の GetButtonState で受け取った Pin 番号の比較にも使われる。
const (
	WioKeyA    = 10
	WioKeyB    = 11
	WioKeyC    = 12
	Wio5SUp    = 13
	Wio5SDown  = 14
	Wio5SLeft  = 15
	Wio5SRight = 16
	Wio5SPress = 17
)
