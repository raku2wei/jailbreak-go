package console_test

import (
	"strings"
	"testing"

	"jailbreak/assets"
	"jailbreak/internal/console"
	"jailbreak/internal/room"
)

// rowText は VTerm の row 行目を文字列にする(全角の継続マーク R=0 は除く)。
func rowText(cells []console.Cell, cols, row int) string {
	var b strings.Builder
	for x := 0; x < cols; x++ {
		if r := cells[row*cols+x].R; r != 0 {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// colOf は row 行目で target が最初に現れるセル位置を返す(見つからなければ -1)。
func colOf(cells []console.Cell, cols, row int, target rune) int {
	for x := 0; x < cols; x++ {
		if cells[row*cols+x].R == target {
			return x
		}
	}
	return -1
}

// 文字幅はロケール(LANG 等)に依存せず、同梱フォントのグリフ幅に合わせた固定値になること。
func TestVTermRuneWidth(t *testing.T) {
	cases := []struct {
		s    string
		want int // s の直後に書いた "|" の桁
	}{
		{"A", 1},
		{"あ", 2},
		{"＼", 2},
		{"★", 2}, // U+2605: フォント上は全角幅
		{"☆", 2}, // U+2606
		{"↑", 1}, // 矢印は半角幅
		{"→", 1},
		{"─", 1}, // 罫線は半角幅
		{"│", 1},
		{"■", 1},
		{"ω", 1},
	}
	for _, c := range cases {
		v := console.NewVTerm(80, 5)
		v.WriteString(c.s + "|")
		cols, _ := v.Size()
		if got := colOf(v.Snapshot(nil), cols, 0, '|'); got != c.want {
			t.Errorf("%q の幅: got %d, want %d", c.s, got, c.want)
		}
	}
}

// タイトルロゴが80桁に収まり、折り返さないこと。
func TestVTermTitleLogoFits80(t *testing.T) {
	logo := strings.TrimRight(string(assets.MustRead("assets/title/logo")), "\n")
	lines := strings.Split(logo, "\n")

	v := console.NewVTerm(80, 55)
	v.WriteString(logo + "\n")
	cols, _ := v.Size()
	cells := v.Snapshot(nil)
	for i, l := range lines {
		// 折り返しが起きると以降の行がずれて一致しなくなる
		got := strings.TrimRight(rowText(cells, cols, i), " ")
		want := strings.TrimRight(l, " ")
		if got != want {
			t.Errorf("logo %d 行目が80桁に収まっていない\n got: %q\nwant: %q", i+1, got, want)
		}
	}
}

// ヒント部屋の★行で右の壁「／」が通常の部屋と同じ桁に来ること。
func TestVTermHintRoomWallAligned(t *testing.T) {
	wallCol := func(hasHint bool) int {
		v := console.NewVTerm(80, 55)
		console.SetBackend(v)
		r := room.Room{HasHint: hasHint}
		r.Display(room.North)
		cols, _ := v.Size()
		return colOf(v.Snapshot(nil), cols, 1, '／')
	}
	normal, hint := wallCol(false), wallCol(true)
	if normal < 0 || hint != normal {
		t.Errorf("ヒント部屋の壁がずれている: hint=%d, normal=%d", hint, normal)
	}
}
