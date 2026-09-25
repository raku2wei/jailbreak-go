package console_test

import (
	"strings"
	"testing"

	"jailbreak/internal/console"
)

// pos は s を書いた後、マーカー "|" がどこに書かれるかでカーソル位置を調べる。
func pos(t *testing.T, v *console.VTerm, s string) (x, y int) {
	t.Helper()
	v.WriteString(s + "|")
	cols, rows := v.Size()
	cells := v.Snapshot(nil)
	for y := 0; y < rows; y++ {
		if x := colOf(cells, cols, y, '|'); x >= 0 {
			return x, y
		}
	}
	t.Fatalf("マーカーが見つからない")
	return -1, -1
}

func TestVTermCursor(t *testing.T) {
	cases := []struct {
		name         string
		cols, rows   int
		in           string
		wantX, wantY int
	}{
		{"初期位置", 10, 4, "", 0, 0},
		{"文字で右へ", 10, 4, "abc", 3, 0},
		{"\\n は改行+行頭(LF で CR も兼ねる)", 10, 4, "abc\n", 0, 1},
		{"\\r は行頭のみ", 10, 4, "abc\r", 0, 0},
		{"\\r\\n", 10, 4, "abc\r\n", 0, 1},
		{"\\n\\n", 10, 4, "a\n\n", 0, 2},
		{"ちょうど右端まで書いても即座には折り返さない(次の文字で折り返す)", 10, 4, "0123456789", 0, 1},
		{"右端ちょうど+\\n は1行しか進まない", 10, 4, "0123456789\n", 0, 1},
		{"右端を超えると次の行へ折り返す", 10, 4, "0123456789ab", 2, 1},
		{"タブは8桁区切り", 20, 4, "a\t", 8, 0},
		{"タブ(8桁目から)", 20, 4, "01234567\t", 16, 0},
		{"タブで右端に達したら改行", 10, 4, "abc\t\t", 0, 1},
		{"全角は2桁進む", 10, 4, "あい", 4, 0},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(c.cols, c.rows)
			x, y := pos(t, v, c.in)
			if x != c.wantX || y != c.wantY {
				t.Errorf("カーソル: got (%d,%d), want (%d,%d)", x, y, c.wantX, c.wantY)
			}
		})
	}
}

// 最終行で改行すると1行スクロールし、新しい最終行はデフォルト属性の空白で埋まること。
func TestVTermScroll(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string // 各行(右側の空白は除く)
	}{
		{"スクロールなし", "a\nb\nc", []string{"a", "b", "c"}},
		{"最終行で \\n すると1行スクロール", "a\nb\nc\n", []string{"b", "c", ""}},
		{"2行スクロール", "a\nb\nc\nd\ne", []string{"c", "d", "e"}},
		{"右端の折り返しでもスクロール", "a\nb\nccccdd", []string{"b", "cccc", "dd"}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(4, 3)
			v.WriteString(c.in)
			cells := v.Snapshot(nil)
			for y, w := range c.want {
				if got := strings.TrimRight(rowText(cells, 4, y), " "); got != w {
					t.Errorf("%d 行目: got %q, want %q", y, got, w)
				}
			}
		})
	}
}

func TestVTermScrollClearsNewLineAttributes(t *testing.T) {
	v := console.NewVTerm(4, 2)
	v.WriteString(esc + "[7;41mabcd" + esc + "[0m\n\n")
	cells := v.Snapshot(nil)
	for i := 4; i < 8; i++ {
		if cells[i] != (console.Cell{R: ' ', Fg: console.ColorDefault, Bg: console.ColorDefault}) {
			t.Errorf("スクロールで空いた行のセル %d が初期化されていない: %+v", i, cells[i])
		}
	}
}

// Clear 後は全セルが初期値になり、カーソルは左上に戻ること。
// ただし SGR 属性(色・太字・反転)とエスケープのパース状態は Clear ではリセットされない。
func TestVTermClear(t *testing.T) {
	v := console.NewVTerm(8, 3)
	v.WriteString(esc + "[1;7;32;43mあいう\nxyz\nfoo")
	v.Clear()
	cells := v.Snapshot(nil)
	blank := console.Cell{R: ' ', Fg: console.ColorDefault, Bg: console.ColorDefault}
	for i, c := range cells {
		if c != blank {
			t.Fatalf("Clear 後のセル %d が初期化されていない: %+v", i, c)
		}
	}
	x, y := pos(t, v, "")
	if x != 0 || y != 0 {
		t.Errorf("Clear 後のカーソル: got (%d,%d), want (0,0)", x, y)
	}
	// 属性は引き継がれる(実装どおりの挙動を固定)
	if got := attrOf(v.Snapshot(nil)[0]); got != (attr{Fg: 2, Bg: 3, Bold: true, Reverse: true}) {
		t.Errorf("Clear 後の書き込み属性: got %+v(Clear は SGR をリセットしない想定)", got)
	}
}

// 全角文字(2セル)の配置。
func TestVTermWideChars(t *testing.T) {
	cases := []struct {
		name string
		cols int
		in   string
		row0 []rune // 0行目の先頭からのセル R 値
		row1 []rune // 1行目の先頭からのセル R 値
	}{
		{"全角の直後は継続セル R=0", 6, "あい", []rune{'あ', 0, 'い', 0, ' ', ' '}, nil},
		{"半角と全角の混在", 6, "aあb", []rune{'a', 'あ', 0, 'b', ' ', ' '}, nil},
		{"★ も2セル", 6, "★a", []rune{'★', 0, 'a'}, nil},
		{"全角スペース U+3000 も2セル占有", 6, "a　b", []rune{'a', '　', 0, 'b'}, nil},
		{"行末に全角がちょうど収まる", 6, "abcdあ", []rune{'a', 'b', 'c', 'd', 'あ', 0}, []rune{' '}},
		{"行末に1セルしか残らない全角は次の行へ折り返す(最終列は空白のまま)",
			6, "abcdeあ", []rune{'a', 'b', 'c', 'd', 'e', ' '}, []rune{'あ', 0, ' '}},
		{"全角だけで折り返し", 5, "あいう", []rune{'あ', 0, 'い', 0, ' '}, []rune{'う', 0, ' '}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			v := console.NewVTerm(c.cols, 3)
			v.WriteString(c.in)
			cells := v.Snapshot(nil)
			check := func(row int, want []rune) {
				for x, w := range want {
					if got := cells[row*c.cols+x].R; got != w {
						t.Errorf("(%d,%d): got %q, want %q", x, row, got, w)
					}
				}
			}
			check(0, c.row0)
			check(1, c.row1)
		})
	}
}

// 継続セルは本体と同じ属性を持つこと(背景色・反転が2マスとも塗られるため)。
func TestVTermWideCharContinuationAttr(t *testing.T) {
	v := console.NewVTerm(6, 2)
	v.WriteString(esc + "[1;7;33;44mあ")
	cells := v.Snapshot(nil)
	want := attr{Fg: 3, Bg: 4, Bold: true, Reverse: true}
	for i := 0; i < 2; i++ {
		if got := attrOf(cells[i]); got != want {
			t.Errorf("セル %d: got %+v, want %+v", i, got, want)
		}
	}
}

// 全角の1マス目に半角を上書きすると、2マス目の継続セル(R=0)は残る(実装どおりの挙動を固定)。
// 描画側は R=0 を描かないので、見た目上は半角1文字+空白1マスになる。
func TestVTermOverwriteWideCharLeavesContinuation(t *testing.T) {
	v := console.NewVTerm(6, 2)
	v.WriteString("あ\rx")
	cells := v.Snapshot(nil)
	if cells[0].R != 'x' || cells[1].R != 0 {
		t.Errorf("got [%q %q], want ['x' 0]", cells[0].R, cells[1].R)
	}
}
