# wiodisplay

Wio Terminal のデスクトップエミュレーター。TinyGo コードを実機に書き込まずに、PC 上でディスプレイ出力・ボタン/ジョイスティック入力を確認できます。

## アーキテクチャ

```
[wiodisplay サーバー]  ← Ebitengine で Wio Terminal UI を表示・入力を管理
        ↑ TCP (net/rpc)
[ユーザーコード]       ← このモジュールのラッパーパッケージを import
```

サーバーとユーザーコードはプロセスを分離しています。サーバーを起動したまま、ユーザーコードだけを再起動できるためホットリロードが可能です。

## インストール

```sh
git clone https://github.com/kurakura967/wiodisplay
cd wiodisplay
go install ./cmd/wio-emu/
```

## 使い方

### 1. サーバーを起動

```sh
wio-emu
```

Wio Terminal の外観ウィンドウが開き、TCP `:9812` でリッスンを開始します。

### 2. ユーザーコードを実行

ユーザーコードにはこのモジュールのパッケージを import します。

```go
import (
    "github.com/kurakura967/wiodisplay/initdisplay"
    "github.com/kurakura967/wiodisplay/machine"
)
```

| build tag | 動作 |
|---|---|
| `!tinygo`（デスクトップ） | TCP 経由でサーバーに接続 |
| `tinygo`（実機） | 実際のハードウェアドライバを使用 |

import を変えずに、実機とエミュレーターの両方で動作します。

```sh
# デスクトップ（エミュレーター）
go run ./your-app/

# 実機
tinygo flash -target wioterminal ./your-app/
```

### サンプルコード

```go
package main

import (
    "image/color"
    "time"

    "github.com/kurakura967/wiodisplay/initdisplay"
    "github.com/kurakura967/wiodisplay/machine"
)

func main() {
    display := initdisplay.InitDisplay()
    display.FillScreen(color.RGBA{R: 0, G: 0, B: 128, A: 255})

    for {
        if !machine.WIO_KEY_A.Get() {
            display.FillScreen(color.RGBA{R: 200, G: 0, B: 0, A: 255})
        }
        time.Sleep(16 * time.Millisecond)
    }
}
```

## ホットリロード

サーバーを起動したまま、ユーザーコードのプロセスだけを再起動します。[entr](https://eradman.com/entrproject/) が必要です。

```sh
# ターミナル 1: サーバーを起動
wio-emu

# ターミナル 2: .go ファイルの変更を検知して自動再起動
find . -name '*.go' | entr -r go run ./your-app/
```

または Makefile を使う場合：

```sh
make run                            # サーバー起動
make dev TARGET=./examples/hello/   # ホットリロードで example 実行
```

## キーマッピング

| Wio Terminal | キーボード | マウス |
|---|---|---|
| ボタン A | `Z` | 右ボタン領域クリック |
| ボタン B | `X` | 中央ボタン領域クリック |
| ボタン C | `C` | 左ボタン領域クリック |
| 5Way Up | `↑` | — |
| 5Way Down | `↓` | — |
| 5Way Left | `←` | — |
| 5Way Right | `→` | — |
| 5Way Press | `Enter` | ジョイスティック領域クリック |

## パッケージ構成

```
wiodisplay/
├── cmd/wio-emu/       サーバーエントリポイント
├── server/            TCP サーバー・Ebitengine 統合
├── client/            共有 RPC コネクション・引数型
├── initdisplay/       InitDisplay() — サーバー接続 or 実機初期化
├── driver/ili9341/    描画 API — RPC 経由 or 実機ドライバ
├── machine/           Pin 型・WIO_KEY_* 定数 — RPC 経由 or 実機 machine
└── examples/          動作確認用サンプル
```

## RPC API

サーバーは以下の RPC メソッドを提供します。

| メソッド | 説明 |
|---|---|
| `DisplayService.DrawPixel` | 1 ピクセル描画 |
| `DisplayService.FillRectangle` | 矩形塗りつぶし |
| `DisplayService.FillScreen` | 画面全体塗りつぶし |
| `DisplayService.GetButtonState` | ボタン押下状態取得 |

## 動作環境

- Go 1.23 以上
- macOS / Linux / Windows
- TinyGo（実機ビルド時）
