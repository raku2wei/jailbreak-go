package console

import (
	"fmt"
	"strings"
)

// 画面の列数・行数。Ebitengine 版の仮想端末の格子と、ターミナル版の推奨サイズ(README)に共通の値。
// 画面内の位置決め(中央寄せ等)はこの大きさを前提にする。
const (
	ScreenCols = 80
	ScreenRows = 55
)

// Area は画面上の連続した行の範囲。Top は 0 始まりの行番号で、幅は画面全体(ScreenCols)。
type Area struct {
	Top  int // 先頭の行(0 始まり)
	Rows int // 行数
}

// ダンジョン画面の区画(0 始まりの行番号)。
// 部屋の描画(room.Display)・操作説明(dungeon.WaitAction)・マップ(dungeon.PrintMap)は
// 画面を消したあと上から順に出力しているので、この行割りはそれらの出力内容で決まる。
// 実際の出力と一致していることは internal/dungeon のテストで確認している。
//
//	 0〜22 行目  一人称視点の描写(上の壁2行 + 部屋のAA 21行)
//	23〜24 行目  空行
//	25    行目  操作説明(w：前に進む …)
//	26〜49 行目  マップ(先頭の空行 + 部屋20行 + 方位3行)
var (
	// ViewArea は一人称視点の描写エリア。戦闘画面はこの中に描く。
	ViewArea = Area{Top: 0, Rows: 23}
	// GuideArea は操作説明の行。
	GuideArea = Area{Top: 25, Rows: 1}
	// MapArea はマップの描画エリア。
	MapArea = Area{Top: 26, Rows: 24}
)

// CenterTop は高さ h 行のブロックをエリアの上下中央に置くときの先頭行(画面上の行番号)を返す。
// エリアより高い場合はエリアの先頭行。
func (a Area) CenterTop(h int) int {
	return a.Top + max(0, (a.Rows-h)/2)
}

// TextWidth は s を画面に出力したときに占めるセル数を返す。
// ANSI エスケープは幅に数えない。文字幅は VTerm と同じ判定(runeWidth)を使うので、
// 全角は2セル、★/☆ も2セルとして数える。改行は含まない前提(1行分の幅)。
func TextWidth(s string) int {
	w := 0
	inEsc, inCSI := false, false
	for _, r := range s {
		switch {
		case inCSI:
			if r >= 0x40 && r <= 0x7e { // CSI の終端文字
				inCSI = false
			}
		case inEsc:
			// VTerm と同じく、ESC の次の1文字は表示しない('[' なら CSI の開始)
			inEsc = false
			inCSI = r == '['
		case r == 0x1b:
			inEsc = true
		default:
			if n := runeWidth(r); n > 0 {
				w += n
			}
		}
	}
	return w
}

// BlockWidth は lines のうち最も幅の広い行のセル数を返す。
func BlockWidth(lines []string) int {
	w := 0
	for _, l := range lines {
		w = max(w, TextWidth(l))
	}
	return w
}

// CenterLeft は幅 w セルのブロックを画面の左右中央に置くときの左余白(セル数)を返す。
// 画面より広い場合は 0。
func CenterLeft(w int) int {
	return max(0, (ScreenCols-w)/2)
}

// TrimBlankLines は先頭・末尾の空行(空白・全角スペースのみの行を含む)を取り除く。
// AA ファイル前後の余白行を、中央寄せの高さ計算に含めないために使う。
func TrimBlankLines(lines []string) []string {
	blank := func(s string) bool { return strings.TrimSpace(s) == "" }
	for len(lines) > 0 && blank(lines[0]) {
		lines = lines[1:]
	}
	for len(lines) > 0 && blank(lines[len(lines)-1]) {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// Frame は画面の一部を書き換える出力を組み立て、Flush で1回の書き込みとして出力する。
// 位置決めはカーソル移動(CSI 行;列 H)、消去は行消去(CSI 2K)で行うので、
// 画面全体を消さずに指定したエリアだけを描き直せる(VTerm・実ターミナルの両方が解釈する)。
// まとめて1回で書くので、VTerm では消去と描画の途中の状態が画面に出ない。
type Frame struct {
	b strings.Builder
}

// MoveTo はカーソルを row 行 col 列(どちらも 0 始まり)へ移動する。
func (f *Frame) MoveTo(row, col int) {
	fmt.Fprintf(&f.b, "\x1b[%d;%dH", row+1, col+1)
}

// Erase はエリアの全行を消去し、カーソルをエリアの先頭行の行頭に置く。
// 消去の前に SGR をリセットするので、消した行に背景色は残らない。
func (f *Frame) Erase(a Area) {
	f.b.WriteString("\x1b[0m")
	for y := a.Top; y < a.Top+a.Rows; y++ {
		f.MoveTo(y, 0)
		f.b.WriteString("\x1b[2K")
	}
	f.MoveTo(a.Top, 0)
}

// Lines は lines を top 行目から1行ずつ、left 列目を左端として書く。
// 余白はカーソル移動で作るので、背景色付きの行でも余白は塗られない。空行は何も書かない。
func (f *Frame) Lines(top, left int, lines []string) {
	for i, l := range lines {
		if l == "" {
			continue
		}
		f.MoveTo(top+i, left)
		f.b.WriteString(l)
	}
}

// Print は現在のカーソル位置から s をそのまま書く。
func (f *Frame) Print(s string) {
	f.b.WriteString(s)
}

// Flush は組み立てた出力をバックエンドへ書き込み、Frame を空にする。
func (f *Frame) Flush() {
	backend.WriteString(f.b.String())
	f.b.Reset()
}
