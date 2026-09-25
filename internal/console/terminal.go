//go:build !js

package console

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"jailbreak/pkg/system"

	"github.com/mattn/go-tty"
)

// TerminalBackend は従来どおり実ターミナルに入出力するバックエンド。
// 出力: 標準出力へそのまま書き込む(ANSIエスケープはターミナルが解釈する)
// 入力: go-tty で読み取ったキーをチャネル経由で提供する
//
// go-tty は Open 時に端末の ECHO/ICANON を落とし、Close で元の端末モードに戻す。
// tty はプロセス生存中開いたままにするので、終了時には必ず Close を呼ぶこと
// (呼ばないとシェルに戻ったあとも入力が表示されない状態が残る)。
type TerminalBackend struct {
	input     chan rune
	tt        *tty.TTY
	closeOnce sync.Once
	closed    atomic.Bool
}

// NewTerminalBackend は実ターミナル用バックエンドを生成する。
// プロセス生存中は tty を開いたままにし、入力読み取り用 goroutine を常駐させる。
func NewTerminalBackend() *TerminalBackend {
	tt, err := tty.Open()
	if err != nil {
		// まだ端末モードを変更していないので、そのまま終了してよい
		log.Fatal(err)
	}
	t := &TerminalBackend{input: make(chan rune, 64), tt: tt}

	go func() {
		for {
			r, err := tt.ReadRune()
			if err != nil {
				if t.closed.Load() {
					return
				}
				t.Fatal(err)
			}
			t.input <- r
		}
	}()

	return t
}

// Close は tty を閉じ、端末モード(echo/icanon 等)を Open 前の状態に戻す。
// 複数回・複数 goroutine から呼んでも安全。
func (t *TerminalBackend) Close() error {
	var err error
	t.closeOnce.Do(func() {
		t.closed.Store(true)
		err = t.tt.Close()
	})
	return err
}

// Fatal は端末モードを戻してからエラーを出力して終了する(log.Fatal の代わり)。
// log.Fatal は defer を実行せずに終了するため、直接呼ぶと端末モードが戻らない。
func (t *TerminalBackend) Fatal(v ...interface{}) {
	t.Close()
	log.Fatal(v...)
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
