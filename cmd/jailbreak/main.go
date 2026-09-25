//go:build !js

// 従来どおりUNIXターミナル上で動作するエントリポイント。
package main

import (
	"os"
	"os/signal"
	"syscall"

	"jailbreak/internal/console"
	"jailbreak/internal/game"
)

func main() {
	// 実ターミナル用バックエンド(標準出力 + go-tty)を設定
	tb := console.NewTerminalBackend()
	// 正常終了(タイトルで「ゲーム終了」)・panic 時に端末モードを戻す
	defer tb.Close()

	// Ctrl-C / SIGTERM でも端末モードを戻してから終了する
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		tb.Close()
		// シェルの慣例(128 + シグナル番号)に合わせた終了コードにする
		code := 1
		if s, ok := sig.(syscall.Signal); ok {
			code = 128 + int(s)
		}
		os.Exit(code)
	}()

	console.SetBackend(tb)
	game.Run()
}
