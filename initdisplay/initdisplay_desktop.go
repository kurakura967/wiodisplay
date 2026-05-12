//go:build !tinygo

package initdisplay

import (
	"log"
	"net/rpc"
	"time"

	"github.com/kurakura967/wiodisplay/client"
	"github.com/kurakura967/wiodisplay/driver/ili9341"
	"github.com/kurakura967/wiodisplay/machine"
)

const serverAddr = "127.0.0.1:9812"

// InitDisplay はサーバーへ接続し、RPC 経由で描画できる Display を返す。
// サーバーが起動するまで最大 5 秒間リトライする。
func InitDisplay() *ili9341.Device {
	var c *rpc.Client
	var err error
	for i := 0; i < 10; i++ {
		c, err = rpc.Dial("tcp", serverAddr)
		if err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if err != nil {
		log.Fatalf("wio-emu: failed to connect to server at %s: %v", serverAddr, err)
	}
	client.Conn = c
	return ili9341.NewSPI(machine.SPI3, machine.LCD_DC, machine.LCD_SS_PIN, machine.LCD_RESET)
}
