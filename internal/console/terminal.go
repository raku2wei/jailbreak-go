//go:build !js

package console

import (
	"fmt"
	"log"

	"jailbreak/pkg/system"

	"github.com/mattn/go-tty"
)

// TerminalBackend は従来どおり実ターミナルに入出力するバックエンド。
// 出力: 標準出力へそのまま書き込む(ANSIエスケープはターミナルが解釈する)
// 入力: go-tty で読み取ったキーをチャネル経由で提供する
type TerminalBackend struct {
	input chan rune
}

// NewTerminalBackend は実ターミナル用バックエンドを生成する。
// プロセス生存中は tty を開いたままにし、入力読み取り用 goroutine を常駐させる。
func NewTerminalBackend() *TerminalBackend {
	t := &TerminalBackend{input: make(chan rune, 64)}

	tt, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			r, err := tt.ReadRune()
			if err != nil {
				log.Fatal(err)
			}
			t.input <- r
		}
	}()

	return t
}

func (t *TerminalBackend) WriteString(s string) {
	fmt.Print(s)
}

func (t *TerminalBackend) Clear() {
	system.System("clear")
}

func (t *TerminalBackend) ReadKey() rune {
	return <-t.input
}

func (t *TerminalBackend) TryReadKey() (rune, bool) {
	select {
	case r := <-t.input:
		return r, true
	default:
		return 0, false
	}
}
