//go:build !js

// WASM版をローカルで動作確認するための簡易HTTPサーバー。
//
// 使い方:
//
//	GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o web/jailbreak.wasm ./cmd/jailbreak-ebiten
//	go run ./cmd/serve
//	→ ブラウザで http://localhost:8080 を開く
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	addr := flag.String("addr", ":8080", "待ち受けアドレス")
	dir := flag.String("dir", "web", "配信するディレクトリ")
	flag.Parse()

	log.Printf("http://localhost%s で %s/ を配信します", *addr, *dir)
	log.Fatal(http.ListenAndServe(*addr, http.FileServer(http.Dir(*dir))))
}
