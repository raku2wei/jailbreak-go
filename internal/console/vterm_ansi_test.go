package console_test

import (
	"testing"

	"jailbreak/internal/console"
)

// VTerm が解釈するエスケープシーケンス(vterm.go の writeRune / applySGR より):
//
//   - CSI ... m (SGR) のみ解釈する。対応パラメータ:
//     0 リセット / 1 太字 / 7 反転 / 22 太字解除 / 27 反転解除 /
//     30-37 前景色 / 39 前景色デフォルト / 40-47 背景色 / 49 背景色デフォルト。
//     パラメータ省略(ESC[m)は 0 と同じ。";" 区切りで複数指定可。上記以外の数値は無視。
//   - SGR 以外の CSI(終端 0x40-0x7e が 'm' 以外。ESC[2J, ESC[H, ESC[?25l 等)は読み捨てる。
//     ESC[2J でも画面はクリアされない(クリアは Clear() が担当)。
//   - CSI 以外のエスケープ(ESC c 等)は ESC の直後の1文字ごと読み捨てる。
//   - パース状態は WriteString の呼び出しをまたいで保持される(分割されたシーケンスも解釈される)。
//   - 制御文字は \n(改行+行頭)、\r(行頭)、\t(8桁タブ)のみ。BS・BEL 等は幅0として無視。

const esc = "\x1b"

// attr は Cell から文字以外の属性だけを取り出して比較するための型。
type attr struct {
	Fg, Bg        int
	Bold, Reverse bool
}

func attrOf(c console.Cell) attr {
	return attr{Fg: c.Fg, Bg: c.Bg, Bold: c.Bold, Reverse: c.Reverse}
}

var defaultAttr = attr{Fg: console.ColorDefault, Bg: console.ColorDefault}

func TestVTermSGR(t *testing.T) {
	cases := []struct {
		name string
		in   string // この後に "X" を書き、X のセル属性を検査する
		want attr
	}{
		{"属性なし", "", defaultAttr},
		{"前景色 30(黒)", esc + "[30m", attr{Fg: 0, Bg: -1}},
		{"前景色 31(赤)", esc + "[31m", attr{Fg: 1, Bg: -1}},
		{"前景色 37(白)", esc + "[37m", attr{Fg: 7, Bg: -1}},
		{"背景色 40", esc + "[40m", attr{Fg: -1, Bg: 0}},
		{"背景色 44", esc + "[44m", attr{Fg: -1, Bg: 4}},
		{"背景色 47", esc + "[47m", attr{Fg: -1, Bg: 7}},
		{"太字 1", esc + "[1m", attr{Fg: -1, Bg: -1, Bold: true}},
		{"反転 7", esc + "[7m", attr{Fg: -1, Bg: -1, Reverse: true}},
		{"複数パラメータ 1;31", esc + "[1;31m", attr{Fg: 1, Bg: -1, Bold: true}},
		{"複数パラメータ 1;7;32;45", esc + "[1;7;32;45m", attr{Fg: 2, Bg: 5, Bold: true, Reverse: true}},
		{"別々に指定しても累積する", esc + "[1m" + esc + "[33m" + esc + "[46m", attr{Fg: 3, Bg: 6, Bold: true}},
		{"後勝ち 31;32", esc + "[31;32m", attr{Fg: 2, Bg: -1}},
		{"リセット 0", esc + "[1;7;31;41m" + esc + "[0m", defaultAttr},
		{"パラメータ省略は 0 扱い", esc + "[1;31m" + esc + "[m", defaultAttr},
		{"1;31 の途中で 0", esc + "[1;31;0;34m", attr{Fg: 4, Bg: -1}},
		{"22 で太字解除", esc + "[1;31m" + esc + "[22m", attr{Fg: 1, Bg: -1}},
		{"27 で反転解除", esc + "[7;41m" + esc + "[27m", attr{Fg: -1, Bg: 1}},
		{"39 で前景色のみデフォルト", esc + "[1;31;42m" + esc + "[39m", attr{Fg: -1, Bg: 2, Bold: true}},
		{"49 で背景色のみデフォルト", esc + "[31;42m" + esc + "[49m", attr{Fg: 1, Bg: -1}},
		{"先頭ゼロ 031", esc + "[031m", attr{Fg: 1, Bg: -1}},
		{"未対応 SGR(4 下線)は無視", esc + "[31m" + esc + "[4m", attr{Fg: 1, Bg: -1}},
		{"256色指定 38;5;196 は無視(色も変わらない)", esc + "[31m" + esc + "[38;5;196m", attr{Fg: 1, Bg: -1}},
		{"明るい色 90-97 は無視", esc + "[31m" + esc + "[91m", attr{Fg: 1, Bg: -1}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(20, 3)
			v.WriteString(c.in + "X")
			cells := v.Snapshot(nil)
			if cells[0].R != 'X' {
				t.Fatalf("エスケープが文字として出力された: 先頭セル=%q, 行=%q", cells[0].R, rowText(cells, 20, 0))
			}
			if got := attrOf(cells[0]); got != c.want {
				t.Errorf("属性: got %+v, want %+v", got, c.want)
			}
		})
	}
}

// 属性は書いた時点の値がセルに固定され、後の SGR で過去のセルは変わらないこと。
func TestVTermSGRAppliesOnlyToLaterCells(t *testing.T) {
	v := console.NewVTerm(20, 3)
	v.WriteString("a" + esc + "[31mb" + esc + "[0mc")
	cells := v.Snapshot(nil)
	want := []attr{defaultAttr, {Fg: 1, Bg: -1}, defaultAttr}
	for i, w := range want {
		if got := attrOf(cells[i]); got != w {
			t.Errorf("セル %d (%q): got %+v, want %+v", i, cells[i].R, got, w)
		}
	}
	if got := rowText(cells, 20, 0)[:3]; got != "abc" {
		t.Errorf("文字列: got %q, want %q", got, "abc")
	}
}

// 未対応のシーケンスが来ても落ちず、文字として出力されず、属性やカーソルが壊れないこと。
func TestVTermUnsupportedSequences(t *testing.T) {
	cases := []struct {
		name string
		seq  string
	}{
		{"カーソル非表示 ?25l", esc + "[?25l"},
		{"カーソル表示 ?25h", esc + "[?25h"},
		{"画面消去 2J", esc + "[2J"},
		{"カーソルホーム H", esc + "[H"},
		{"カーソル位置 10;5H", esc + "[10;5H"},
		{"カーソル上 A", esc + "[3A"},
		{"行消去 K", esc + "[K"},
		{"CSI 以外 ESC c", esc + "c"},
		{"CSI 以外 ESC 7", esc + "7"},
		{"CSI 以外 ESC (B", esc + "(" + "B"}, // "(" だけ読み捨て、"B" は文字として残る
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(20, 3)
			v.WriteString(esc + "[1;31m" + "ab")
			v.WriteString(c.seq)
			v.WriteString("cd")
			cells := v.Snapshot(nil)

			want := "abcd"
			if c.name == "CSI 以外 ESC (B" {
				want = "abBcd"
			}
			got := rowText(cells, 20, 0)[:len(want)]
			if got != want {
				t.Errorf("行: got %q, want %q", got, want)
			}
			// 画面は消えておらず、カーソルも移動していない(直前の続きに書かれる)
			for i := range want {
				if a := attrOf(cells[i]); a != (attr{Fg: 1, Bg: -1, Bold: true}) {
					t.Errorf("セル %d の属性が壊れた: %+v", i, a)
				}
			}
			if r := rowText(cells, 20, 1); r != "                    " {
				t.Errorf("2行目に何か書かれた: %q", r)
			}
		})
	}
}

// ESC シーケンスが WriteString をまたいで分割されても解釈されること。
func TestVTermSplitEscapeAcrossWrites(t *testing.T) {
	v := console.NewVTerm(20, 3)
	v.WriteString("a" + esc)
	v.WriteString("[3")
	v.WriteString("2mb")
	cells := v.Snapshot(nil)
	if got := rowText(cells, 20, 0)[:2]; got != "ab" {
		t.Fatalf("行: got %q, want %q", got, "ab")
	}
	if got := attrOf(cells[1]); got != (attr{Fg: 2, Bg: -1}) {
		t.Errorf("b の属性: got %+v", got)
	}
}

// CSI の終端(0x40-0x7e)が来るまでは後続の文字がすべてシーケンスとして飲み込まれる。
// 例えば "ESC[" の直後に全角文字が続くと、次の英字までが表示されない。
func TestVTermUnterminatedCSISwallowsText(t *testing.T) {
	v := console.NewVTerm(20, 3)
	v.WriteString(esc + "[" + "12あい" + "zXY")
	cells := v.Snapshot(nil)
	// "12あい" は中間バイト扱い、'z' が終端(SGR 以外なので無視)、"XY" から表示
	if got := rowText(cells, 20, 0)[:2]; got != "XY" {
		t.Errorf("行: got %q, want %q", got, "XY")
	}
}

// 表示されない制御文字(BS, BEL 等)は幅0として無視され、カーソルも動かないこと。
func TestVTermIgnoredControlChars(t *testing.T) {
	cases := []struct {
		name string
		ch   string
	}{
		{"BS", "\b"},
		{"BEL", "\a"},
		{"NUL", "\x00"},
		{"DEL", "\x7f"},
		{"結合文字 U+0301", "́"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(10, 2)
			v.WriteString("ab" + c.ch + "c")
			if got := rowText(v.Snapshot(nil), 10, 0)[:3]; got != "abc" {
				t.Errorf("行: got %q, want %q", got, "abc")
			}
		})
	}
}
