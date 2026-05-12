.PHONY: build install run dev

# サーバーバイナリをビルド
build:
	go build -o wio-emu ./cmd/wio-emu/

# サーバーバイナリを $GOPATH/bin にインストール
install:
	go install ./cmd/wio-emu/

# サーバーを起動（別ターミナルで実行）
run:
	./wio-emu

# ホットリロードで example を実行（entr が必要）
# 使い方: make dev TARGET=./examples/hello/
dev:
	find . -name '*.go' -not -path './cmd/*' | entr -r go run $(TARGET)
