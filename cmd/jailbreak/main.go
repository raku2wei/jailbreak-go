//go:build !js

// 従来どおりUNIXターミナル上で動作するエントリポイント。
package main

import (
	"jailbreak/internal/console"
	"jailbreak/internal/game"
)

func main() {
	// 実ターミナル用バックエンド(標準出力 + go-tty)を設定
	console.SetBackend(console.NewTerminalBackend())
	game.Run()
}
