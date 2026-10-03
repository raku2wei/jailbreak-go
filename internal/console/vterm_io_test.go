package console_test

import (
	"sync"
	"testing"
	"time"

	"jailbreak/internal/console"
)

// 入力チャネルの容量(vterm.go の make(chan rune, 64))。
const inputCap = 64

// VTerm は console.Backend を満たすこと。
var _ console.Backend = (*console.VTerm)(nil)

// PushRune はブロックせず、容量を超えた分は捨てる。取り出しは FIFO。
func TestVTermInputOverflowDrops(t *testing.T) {
	cases := []struct {
		name  string
		push  int
		wantN int
	}{
		{"空", 0, 0},
		{"1件", 1, 1},
		{"容量ちょうど", inputCap, inputCap},
		{"容量+1 は1件捨てる", inputCap + 1, inputCap},
		{"容量の2倍", inputCap * 2, inputCap},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(10, 2)
			done := make(chan struct{})
			go func() {
				for i := 0; i < c.push; i++ {
					v.PushRune(rune('A' + i%26))
				}
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("PushRune がブロックした")
			}
			// 先に入れた順に残り、あふれた後半が捨てられる
			for i := 0; i < c.wantN; i++ {
				r, ok := v.TryReadKey()
				if !ok {
					t.Fatalf("%d 件目が読めない(want %d 件)", i, c.wantN)
				}
				if want := rune('A' + i%26); r != want {
					t.Fatalf("%d 件目: got %q, want %q", i, r, want)
				}
			}
			if r, ok := v.TryReadKey(); ok {
				t.Errorf("余分な入力が残っている: %q", r)
			}
		})
	}
}

// TryReadKey は空のときブロックせず (0, false) を返す。
func TestVTermTryReadKeyEmpty(t *testing.T) {
	v := console.NewVTerm(10, 2)
	r, ok := v.TryReadKey()
	if ok || r != 0 {
		t.Errorf("got (%q, %v), want (0, false)", r, ok)
	}
}

// ReadKey は入力が来るまでブロックし、来たら返る。
func TestVTermReadKeyBlocksUntilPush(t *testing.T) {
	v := console.NewVTerm(10, 2)
	got := make(chan rune, 1)
	go func() { got <- v.ReadKey() }()

	select {
	case r := <-got:
		t.Fatalf("入力前に ReadKey が返った: %q", r)
	case <-time.After(20 * time.Millisecond):
	}
	v.PushRune('あ')
	select {
	case r := <-got:
		if r != 'あ' {
			t.Errorf("got %q, want %q", r, 'あ')
		}
	case <-time.After(time.Second):
		t.Fatal("PushRune 後も ReadKey が返らない")
	}
}

// ReadKey はバッファ済みの入力を FIFO で即座に返す。
func TestVTermReadKeyFIFO(t *testing.T) {
	v := console.NewVTerm(10, 2)
	for _, r := range "wasd" {
		v.PushRune(r)
	}
	for _, want := range "wasd" {
		if got := v.ReadKey(); got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	}
}

// 世代番号は Clear / WriteString のたびに1ずつ増える。
// 内容が変わらない書き込み・空文字列・エスケープのみでも増える(差分判定はしない)。
func TestVTermGen(t *testing.T) {
	cases := []struct {
		name string
		op   func(v *console.VTerm)
		want int64 // 操作による増分
	}{
		{"WriteString", func(v *console.VTerm) { v.WriteString("a") }, 1},
		{"WriteString 空文字列", func(v *console.VTerm) { v.WriteString("") }, 1},
		{"WriteString エスケープのみ", func(v *console.VTerm) { v.WriteString(esc + "[31m") }, 1},
		{"長い文字列でも1回の呼び出しで +1", func(v *console.VTerm) { v.WriteString("abc\ndef\nあいう") }, 1},
		{"同じ内容を上書き(\\r で戻って同じ文字)", func(v *console.VTerm) { v.WriteString("\r ") }, 1},
		{"Clear", func(v *console.VTerm) { v.Clear() }, 1},
		{"すでに空の画面を Clear", func(v *console.VTerm) { v.Clear(); v.Clear() }, 2},
		{"Snapshot では増えない", func(v *console.VTerm) { v.Snapshot(nil) }, 0},
		{"PushRune / TryReadKey では増えない", func(v *console.VTerm) { v.PushRune('x'); v.TryReadKey() }, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(10, 3)
			before := v.Gen()
			c.op(v)
			if got := v.Gen() - before; got != c.want {
				t.Errorf("gen の増分: got %d, want %d", got, c.want)
			}
		})
	}
}

func TestVTermGenInitialZero(t *testing.T) {
	if g := console.NewVTerm(10, 3).Gen(); g != 0 {
		t.Errorf("初期 gen: got %d, want 0", g)
	}
}

// Snapshot は buf が足りれば再利用し、足りなければ新規確保する。返り値はコピーで、以後の書き込みの影響を受けない。
func TestVTermSnapshotBuffer(t *testing.T) {
	v := console.NewVTerm(4, 2)
	v.WriteString("ab")

	buf := make([]console.Cell, 8)
	got := v.Snapshot(buf)
	if &got[0] != &buf[0] {
		t.Error("十分な長さの buf が再利用されていない")
	}
	short := make([]console.Cell, 3)
	if got2 := v.Snapshot(short); len(got2) != 8 || &got2[0] == &short[0] {
		t.Errorf("短い buf で新規確保されていない: len=%d", len(got2))
	}
	v.WriteString("\rxy")
	if got[0].R != 'a' {
		t.Errorf("Snapshot の結果が後の書き込みで変わった: %q", got[0].R)
	}
}

// 書き込み(ゲームロジック goroutine)と読み取り(Ebitengine の Update/Draw)を並行に回しても
// データ競合が起きないこと。go test -race で検出する。
func TestVTermConcurrentAccess(t *testing.T) {
	v := console.NewVTerm(80, 24)
	stop := make(chan struct{})
	var wg sync.WaitGroup

	// 書き込み側: 出力と Clear とキー読み取り
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; ; i++ {
			select {
			case <-stop:
				return
			default:
			}
			v.WriteString(esc + "[1;31mあいう" + esc + "[0m abc\n")
			if i%50 == 0 {
				v.Clear()
			}
			v.TryReadKey()
		}
	}()
	// 読み取り側: Draw 相当(Gen を見て変化があれば Snapshot)とキー入力
	wg.Add(1)
	go func() {
		defer wg.Done()
		var buf []console.Cell
		var last int64 = -1
		for {
			select {
			case <-stop:
				return
			default:
			}
			if g := v.Gen(); g != last {
				buf = v.Snapshot(buf)
				last = g
			}
			v.Size()
			v.PushRune('k')
		}
	}()

	time.Sleep(100 * time.Millisecond)
	close(stop)
	wg.Wait()
}
