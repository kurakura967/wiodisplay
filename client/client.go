//go:build !tinygo

package client

import "net/rpc"

// Conn はサーバーへの共有 RPC コネクション。
// initdisplay.InitDisplay() が初期化する。
var Conn *rpc.Client

// RPC メソッドの引数型（server/server.go の型と field 名・型が一致している必要がある）

type DrawPixelArgs struct {
	X, Y  int16
	Color uint16
}

type FillRectangleArgs struct {
	X, Y, W, H int16
	Color      uint16
}

type FillScreenArgs struct {
	Color uint16
}

type GetButtonStateArgs struct {
	Pin int
}
