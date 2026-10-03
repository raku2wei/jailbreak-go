// Package console はゲームロジックとターミナル入出力を分離するための抽象化層。
//
// ゲームロジックは fmt.Printf / go-tty を直接使う代わりに本パッケージを経由する。
// バックエンドとして以下の2種類を提供する:
//   - TerminalBackend: 従来どおり標準出力 + go-tty で実ターミナルに入出力する (terminal.go / js以外)
//   - VTerm: ANSIエスケープをパースして文字グリッドに反映する仮想端末 (vterm.go / Ebitengine・WASM用)
package console

import (
	"fmt"

	"jailbreak/assets"
)

// Backend はゲームの入出力先(実ターミナル or 仮想スクリーン)の抽象。
type Backend interface {
	// WriteString はANSIエスケープを含み得る文字列を出力する。
	WriteString(s string)
	// Clear は画面をクリアする。
	Clear()
	// ReadKey はキー入力1文字を待つ(ブロッキング)。
	ReadKey() rune
	// TryReadKey は入力バッファにキー入力があれば取り出す(ノンブロッキング)。
	TryReadKey() (rune, bool)
}

var backend Backend

// SetBackend は入出力バックエンドを設定する。
// 各エントリポイント(cmd/jailbreak, cmd/jailbreak-ebiten)が起動時に必ず呼ぶこと。
func SetBackend(b Backend) {
	backend = b
}

// Printf は fmt.Printf 相当の出力を行う。
func Printf(format string, a ...interface{}) {
	backend.WriteString(fmt.Sprintf(format, a...))
}

// Println は fmt.Println 相当の出力を行う。
func Println(a ...interface{}) {
	backend.WriteString(fmt.Sprintln(a...))
}

// Clear は画面をクリアする(従来の system.System("clear") 相当)。
func Clear() {
	backend.Clear()
}

// PrintFile は埋め込みアセットのテキストを表示する(従来の system.PrintFile 相当)。
func PrintFile(path string) {
	// 従来実装は fmt.Println(string(bytes)) だったので末尾に改行を付ける
	backend.WriteString(string(assets.MustRead(path)) + "\n")
}

// ReadKey はキー入力1文字を待つ。
func ReadKey() rune {
	return backend.ReadKey()
}

// TryReadKey は入力があれば取り出す(ノンブロッキング)。
func TryReadKey() (rune, bool) {
	return backend.TryReadKey()
}
