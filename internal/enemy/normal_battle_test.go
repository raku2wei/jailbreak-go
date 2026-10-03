package enemy

import (
	"strings"
	"testing"
	"time"

	"jailbreak/assets"
	"jailbreak/internal/console"
)

// screenRows は VTerm の各行を文字列にする(全角の継続マーク R=0 は除き、右側の空白は落とす)。
func screenRows(v *console.VTerm) []string {
	cols, rows := v.Size()
	cells := v.Snapshot(nil)
	out := make([]string, rows)
	for y := 0; y < rows; y++ {
		var b strings.Builder
		for x := 0; x < cols; x++ {
			if r := cells[y*cols+x].R; r != 0 {
				b.WriteRune(r)
			}
		}
		out[y] = strings.TrimRight(b.String(), " ")
	}
	return out
}

func testEnemy(jp, romaji string) *Enemy {
	return &Enemy{Name: "警備員", FilePath: "assets/enemy/keibi", TextJapanese: jp, TextRomaji: romaji}
}

// 戦闘画面(時間切れの最終フレーム)が 80x55 の画面の中央に描かれること。
// AA(幅11)は左余白34、問題文以下(最も広い「残り時間」の行が幅19)は左余白30、
// 全体の高さ15行(AA5 + 空行3 + 問題文以下7)を上下中央に置くので上余白20。
func TestRefreshCentered(t *testing.T) {
	v := console.NewVTerm(console.ScreenCols, console.ScreenRows)
	console.SetBackend(v)

	e := testEnemy("一刀両断", "ittouryoudann")
	if ok := refresh(e, newBattleLayout(e), false, "ittou", 0, time.Now().Add(-time.Minute)); ok {
		t.Fatal("時間切れなのに refresh が true を返した")
	}

	pad := func(n int, s string) string { return strings.Repeat(" ", n) + s }
	want := map[int]string{
		19: "",
		20: pad(34, "  　∧　∧"),
		21: pad(34, " 　( ･ω・)"),
		22: pad(34, "  ⊂      ⊃"),
		23: pad(34, "   ｜    ｜"),
		24: pad(34, "    ∪￣∪"),
		25: "",
		27: "",
		28: pad(30, "「一刀両断」"),
		29: pad(30, "ittouryoudann"),
		30: pad(30, "ittou"),
		31: "",
		32: pad(30, "残り時間 : 0.00 秒"),
		33: pad(30, "TIME OVER!"),
		34: "",
	}
	got := screenRows(v)
	for y, w := range want {
		if got[y] != w {
			t.Errorf("%d 行目: got %q, want %q", y, got[y], w)
		}
	}
	for y, s := range got {
		if _, ok := want[y]; !ok && (y < 19 || y > 34) && s != "" {
			t.Errorf("%d 行目に余計な表示: %q", y, s)
		}
	}
}

// ミス表示の赤背景が左余白に付かないこと。
func TestRefreshMissBackgroundNotOnPadding(t *testing.T) {
	v := console.NewVTerm(console.ScreenCols, console.ScreenRows)
	console.SetBackend(v)

	e := testEnemy("一刀両断", "ittouryoudann")
	l := newBattleLayout(e)
	refresh(e, l, true, "itt", 'x', time.Now())

	cols, _ := v.Size()
	cells := v.Snapshot(nil)
	inputRow := l.top + len(l.aa) + battleGapRows + 2
	for _, y := range []int{inputRow, inputRow + 1} {
		for x := 0; x < l.textLeft; x++ {
			if c := cells[y*cols+x]; c.Bg != console.ColorDefault {
				t.Errorf("(%d,%d) の余白に背景色 %d が付いている", x, y, c.Bg)
			}
		}
	}
	if c := cells[inputRow*cols+l.textLeft+3]; c.R != 'x' || c.Bg != 1 {
		t.Errorf("ミスした文字: got %q bg=%d, want 'x' bg=1(赤)", c.R, c.Bg)
	}
}

// どの問題文でも、問題文以下のブロックが画面からはみ出さず左右中央(余白の差が1以内)に来ること。
func TestBattleLayoutCenteredForAllTexts(t *testing.T) {
	lines := strings.Split(strings.TrimRight(string(assets.MustRead(typingTextFilePath)), "\n"), "\n")
	for i := 0; i+1 < len(lines); i += 2 {
		e := testEnemy(lines[i], lines[i+1])
		l := newBattleLayout(e)
		w := console.TextWidth("「" + e.TextJapanese + "」")
		w = max(w, console.TextWidth(e.TextRomaji), console.TextWidth("残り時間 : 15.00 秒"))
		right := console.ScreenCols - l.textLeft - w
		if right < 0 || l.textLeft-right > 1 || right-l.textLeft > 1 {
			t.Errorf("%q: 左余白 %d / 右余白 %d が中央寄せになっていない", e.TextRomaji, l.textLeft, right)
		}
	}
}
