package server

import (
	"bytes"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/kurakura967/wiodisplay/assets"
)

const (
	ScreenWidth  = 320
	ScreenHeight = 240

	imgWidth  = 1190
	imgHeight = 950

	lcdX = 29
	lcdY = 41
	lcdW = 1133
	lcdH = 801

	RPCPort = ":9812"
)

// Buttons holds the current pressed state of emulated buttons.
// true = pressed. machine.Pin.Get() inverts this (active-low).
var Buttons struct {
	mu                            sync.RWMutex
	A, B, C                       bool
	Up, Down, Left, Right, Center bool
}

// screen is the shared drawing buffer between RPC handlers and Ebitengine.
type screen struct {
	mu     sync.Mutex
	buf    [ScreenWidth * ScreenHeight]color.RGBA
	pixels [ScreenWidth * ScreenHeight * 4]byte
}

var globalScreen = newScreen()

func newScreen() *screen {
	s := &screen{}
	for i := range s.buf {
		s.buf[i] = color.RGBA{0, 0, 0, 255}
	}
	return s
}

func (s *screen) drawPixel(x, y int16, c color.RGBA) {
	if x < 0 || x >= ScreenWidth || y < 0 || y >= ScreenHeight {
		return
	}
	s.mu.Lock()
	s.buf[int(y)*ScreenWidth+int(x)] = c
	s.mu.Unlock()
}

func (s *screen) fillRectangle(x, y, w, h int16, c color.RGBA) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for dy := int16(0); dy < h; dy++ {
		for dx := int16(0); dx < w; dx++ {
			px, py := x+dx, y+dy
			if px < 0 || px >= ScreenWidth || py < 0 || py >= ScreenHeight {
				continue
			}
			s.buf[int(py)*ScreenWidth+int(px)] = c
		}
	}
}

func (s *screen) copyToImage(img *ebiten.Image) {
	s.mu.Lock()
	for i, c := range s.buf {
		s.pixels[i*4] = c.R
		s.pixels[i*4+1] = c.G
		s.pixels[i*4+2] = c.B
		s.pixels[i*4+3] = c.A
	}
	s.mu.Unlock()
	img.WritePixels(s.pixels[:])
}

// RPC argument types

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

// DisplayService は描画・入力 RPC のハンドラ
type DisplayService struct{}

func rgb565ToRGBA(c uint16) color.RGBA {
	r := uint8((c >> 11) & 0x1F)
	g := uint8((c >> 5) & 0x3F)
	b := uint8(c & 0x1F)
	return color.RGBA{
		R: r<<3 | r>>2,
		G: g<<2 | g>>4,
		B: b<<3 | b>>2,
		A: 255,
	}
}

func (d *DisplayService) DrawPixel(args *DrawPixelArgs, reply *struct{}) error {
	globalScreen.drawPixel(args.X, args.Y, rgb565ToRGBA(args.Color))
	return nil
}

func (d *DisplayService) FillRectangle(args *FillRectangleArgs, reply *struct{}) error {
	globalScreen.fillRectangle(args.X, args.Y, args.W, args.H, rgb565ToRGBA(args.Color))
	return nil
}

func (d *DisplayService) FillScreen(args *FillScreenArgs, reply *struct{}) error {
	globalScreen.fillRectangle(0, 0, ScreenWidth, ScreenHeight, rgb565ToRGBA(args.Color))
	return nil
}

func (d *DisplayService) GetButtonState(args *GetButtonStateArgs, reply *bool) error {
	Buttons.mu.RLock()
	defer Buttons.mu.RUnlock()
	switch args.Pin {
	case PinWioKeyA:
		*reply = Buttons.A
	case PinWioKeyB:
		*reply = Buttons.B
	case PinWioKeyC:
		*reply = Buttons.C
	case PinWio5SUp:
		*reply = Buttons.Up
	case PinWio5SDown:
		*reply = Buttons.Down
	case PinWio5SLeft:
		*reply = Buttons.Left
	case PinWio5SRight:
		*reply = Buttons.Right
	case PinWio5SPress:
		*reply = Buttons.Center
	default:
		*reply = false
	}
	return nil
}

// Pin 番号は machine パッケージと合わせる
const (
	PinWioKeyA   = 10
	PinWioKeyB   = 11
	PinWioKeyC   = 12
	PinWio5SUp   = 13
	PinWio5SDown = 14
	PinWio5SLeft = 15
	PinWio5SRight = 16
	PinWio5SPress = 17
)

// StartTCPServer は別ゴルーチンで TCP RPC サーバーを起動する
func StartTCPServer() {
	svc := &DisplayService{}
	if err := rpc.Register(svc); err != nil {
		log.Fatalf("rpc.Register: %v", err)
	}

	ln, err := net.Listen("tcp", RPCPort)
	if err != nil {
		log.Fatalf("listen %s: %v", RPCPort, err)
	}
	log.Printf("wio-emu server listening on %s", RPCPort)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("accept error: %v", err)
				continue
			}
			go rpc.ServeConn(conn)
		}
	}()
}

// Button hit areas in game coordinates (1190x950 space).
type rectHit struct{ x0, y0, x1, y1 int }

func (r rectHit) contains(mx, my int) bool {
	return mx >= r.x0 && mx <= r.x1 && my >= r.y0 && my <= r.y1
}

type circleHit struct{ cx, cy, r int }

func (c circleHit) contains(mx, my int) bool {
	dx, dy := mx-c.cx, my-c.cy
	return dx*dx+dy*dy <= c.r*c.r
}

var (
	hitA        = rectHit{168, 0, 274, 30}
	hitB        = rectHit{389, 0, 495, 30}
	hitC        = rectHit{610, 0, 715, 30}
	hitJoystick = circleHit{1042, 843, 40}
)

// game implements ebiten.Game
type game struct {
	lcdImage    *ebiten.Image
	deviceImage *ebiten.Image
	rawImage    image.Image
	tcpStarted  bool
}

func (g *game) Update() error {
	if !g.tcpStarted {
		StartTCPServer()
		g.tcpStarted = true
	}

	mx, my := ebiten.CursorPosition()
	lmb := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	Buttons.mu.Lock()
	// 実機の並び: 左から C / B / A
	Buttons.C = (lmb && hitA.contains(mx, my)) || ebiten.IsKeyPressed(ebiten.KeyC)
	Buttons.B = (lmb && hitB.contains(mx, my)) || ebiten.IsKeyPressed(ebiten.KeyX)
	Buttons.A = (lmb && hitC.contains(mx, my)) || ebiten.IsKeyPressed(ebiten.KeyZ)
	Buttons.Center = (lmb && hitJoystick.contains(mx, my)) || ebiten.IsKeyPressed(ebiten.KeyEnter)
	Buttons.Up = ebiten.IsKeyPressed(ebiten.KeyArrowUp)
	Buttons.Down = ebiten.IsKeyPressed(ebiten.KeyArrowDown)
	Buttons.Left = ebiten.IsKeyPressed(ebiten.KeyArrowLeft)
	Buttons.Right = ebiten.IsKeyPressed(ebiten.KeyArrowRight)
	Buttons.mu.Unlock()

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	// Metalコンテキストが確立された後（初回Draw時）にebiten.Imageを生成する
	if g.lcdImage == nil {
		g.lcdImage = ebiten.NewImage(ScreenWidth, ScreenHeight)
		g.deviceImage = ebiten.NewImageFromImage(g.rawImage)
	}

	vector.DrawFilledRect(screen, float32(lcdX), float32(lcdY), float32(lcdW), float32(lcdH),
		color.RGBA{0, 0, 0, 255}, false)

	scaleX := float64(lcdW) / float64(ScreenWidth)
	scaleY := float64(lcdH) / float64(ScreenHeight)
	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}
	scaledW := float64(ScreenWidth) * scale
	scaledH := float64(ScreenHeight) * scale
	offsetX := float64(lcdX) + (float64(lcdW)-scaledW)/2
	offsetY := float64(lcdY) + (float64(lcdH)-scaledH)/2

	globalScreen.copyToImage(g.lcdImage)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)
	op.GeoM.Translate(offsetX, offsetY)
	screen.DrawImage(g.lcdImage, op)

	screen.DrawImage(g.deviceImage, nil)

	const pressAlpha = 120
	pressColor := color.RGBA{255, 255, 255, pressAlpha}

	Buttons.mu.RLock()
	a, b, c, center := Buttons.A, Buttons.B, Buttons.C, Buttons.Center
	Buttons.mu.RUnlock()

	if c {
		vector.DrawFilledRect(screen,
			float32(hitA.x0), float32(hitA.y0),
			float32(hitA.x1-hitA.x0), float32(hitA.y1-hitA.y0),
			pressColor, false)
	}
	if b {
		vector.DrawFilledRect(screen,
			float32(hitB.x0), float32(hitB.y0),
			float32(hitB.x1-hitB.x0), float32(hitB.y1-hitB.y0),
			pressColor, false)
	}
	if a {
		vector.DrawFilledRect(screen,
			float32(hitC.x0), float32(hitC.y0),
			float32(hitC.x1-hitC.x0), float32(hitC.y1-hitC.y0),
			pressColor, false)
	}
	if center {
		vector.DrawFilledCircle(screen,
			float32(hitJoystick.cx), float32(hitJoystick.cy),
			float32(hitJoystick.r), pressColor, false)
	}
}

func (g *game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return imgWidth, imgHeight
}

// Run は TCP サーバーを起動してから Ebitengine ウィンドウを開く。
// メインゴルーチンから呼び出すこと。
func Run() {
	img, _, err := image.Decode(bytes.NewReader(assets.WioTerminalBody))
	if err != nil {
		log.Fatalf("loading device image: %v", err)
	}

	ebiten.SetWindowSize(imgWidth/2, imgHeight/2)
	ebiten.SetWindowTitle("Wio Terminal Emulator")

	if err := ebiten.RunGame(&game{rawImage: img}); err != nil {
		log.Fatal(err)
	}
}
