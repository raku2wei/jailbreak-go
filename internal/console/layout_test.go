package console_test

import (
	"strings"
	"testing"

	"jailbreak/internal/console"
)

// TextWidth は VTerm と同じ文字幅で数え、ANSI エスケープは数えないこと。
func TestTextWidth(t *testing.T) {
	cases := []struct {
		s    string
		want int
	}{
		{"", 0},
		{"abc", 3},
		{"あい", 4},
		{"★", 2},
		{"→", 1},
		{"\x1b[1;32mab\x1b[0m", 2},
		{"\x1b[41mタイプミス！\x1b[49m", 12},
	}
	for _, c := range cases {
		if got := console.TextWidth(c.s); got != c.want {
			t.Errorf("TextWidth(%q) = %d, want %d", c.s, got, c.want)
		}
		// VTerm に書いたときのカーソル位置と一致すること
		v := console.NewVTerm(80, 2)
		if x, _ := pos(t, v, c.s); x != c.want {
			t.Errorf("VTerm 上の幅(%q) = %d, want %d", c.s, x, c.want)
		}
	}
}

func TestAreaCenterTop(t *testing.T) {
	cases := []struct {
		a    console.Area
		h    int
		want int
	}{
		{console.Area{Top: 0, Rows: 23}, 15, 4},
		{console.Area{Top: 0, Rows: 23}, 3, 10},
		{console.Area{Top: 10, Rows: 5}, 2, 11},
		{console.Area{Top: 10, Rows: 5}, 9, 10}, // エリアより高いときは先頭行
	}
	for _, c := range cases {
		if got := c.a.CenterTop(c.h); got != c.want {
			t.Errorf("%+v.CenterTop(%d) = %d, want %d", c.a, c.h, got, c.want)
		}
	}
}

// 戦闘画面は一人称視点の描写エリアに収まる必要がある(敵のAA5行 + 空行3行 + 問題文以下7行)。
func TestViewAreaFitsBattle(t *testing.T) {
	if console.ViewArea.Top+console.ViewArea.Rows > console.GuideArea.Top ||
		console.GuideArea.Top+console.GuideArea.Rows > console.MapArea.Top ||
		console.MapArea.Top+console.MapArea.Rows > console.ScreenRows {
		t.Errorf("区画が重なっているか画面からはみ出している: view=%+v guide=%+v map=%+v",
			console.ViewArea, console.GuideArea, console.MapArea)
	}
	if console.ViewArea.Rows < 15 {
		t.Errorf("描写エリアが戦闘画面(15行)より低い: %d 行", console.ViewArea.Rows)
	}
}

// Frame.Erase はエリアの行だけを消し、他の行は残すこと。
// Frame.Lines は指定位置に書き、余白(左側)には何も書かない(背景色も付けない)こと。
func TestFrameEraseAndLines(t *testing.T) {
	v := console.NewVTerm(10, 6)
	console.SetBackend(v)
	v.WriteString(strings.Repeat("##########", 6))

	var f console.Frame
	f.Erase(console.Area{Top: 1, Rows: 3})
	f.Lines(2, 4, []string{"\x1b[41mab\x1b[49m", "", "cd"})
	f.Flush()

	cells := v.Snapshot(nil)
	want := []string{
		"##########",
		"          ",
		"    ab    ",
		"          ",
		"####cd####", // エリア外の行は残る(書いた位置だけ上書き)
		"##########",
	}
	for y, w := range want {
		if got := rowText(cells, 10, y); got != w {
			t.Errorf("%d 行目: got %q, want %q", y, got, w)
		}
	}
	for x := 0; x < 10; x++ {
		c := cells[2*10+x]
		wantBg := console.ColorDefault
		if x == 4 || x == 5 {
			wantBg = 1
		}
		if c.Bg != wantBg {
			t.Errorf("(%d,2) の背景色: got %d, want %d", x, c.Bg, wantBg)
		}
	}
}

// Erase は消去前に SGR をリセットするので、直前の背景色が消した行に残らないこと。
func TestFrameEraseResetsAttributes(t *testing.T) {
	v := console.NewVTerm(4, 2)
	console.SetBackend(v)
	v.WriteString("\x1b[1;44mabcd")

	var f console.Frame
	f.Erase(console.Area{Top: 0, Rows: 2})
	f.Print("x")
	f.Flush()

	cells := v.Snapshot(nil)
	if got := attrOf(cells[0]); got != defaultAttr {
		t.Errorf("消去後に書いた文字の属性: got %+v", got)
	}
	for i := 1; i < len(cells); i++ {
		if cells[i] != (console.Cell{R: ' ', Fg: console.ColorDefault, Bg: console.ColorDefault}) {
			t.Errorf("セル %d が初期値になっていない: %+v", i, cells[i])
		}
	}
}

// Flush は1回の書き込みで出力し(世代番号が1つだけ進む)、Frame を空にすること。
func TestFrameFlushOnce(t *testing.T) {
	v := console.NewVTerm(10, 3)
	console.SetBackend(v)

	var f console.Frame
	f.Erase(console.Area{Top: 0, Rows: 3})
	f.Lines(1, 2, []string{"ab", "cd"})
	g := v.Gen()
	f.Flush()
	if d := v.Gen() - g; d != 1 {
		t.Errorf("Flush の書き込み回数: got %d, want 1", d)
	}
	f.Flush() // 空の Frame の Flush は何も描かない
	if got := rowText(v.Snapshot(nil), 10, 1); got != "  ab      " {
		t.Errorf("2回目の Flush で内容が繰り返された: %q", got)
	}
}
