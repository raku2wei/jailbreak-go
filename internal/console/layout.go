package console

import "strings"

// 画面の列数・行数。Ebitengine 版の仮想端末の格子と、ターミナル版の推奨サイズ(README)に共通の値。
// 画面内の位置決め(中央寄せ等)はこの大きさを前提にする。
const (
	ScreenCols = 80
	ScreenRows = 55
)

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

// CenterTop は高さ h 行のブロックを画面の上下中央に置くときの上余白(行数)を返す。
// 画面より高い場合は 0。
func CenterTop(h int) int {
	return max(0, (ScreenRows-h)/2)
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

// PrintLines は lines を1行ずつ、左に left セル分の空白を付けて出力する。
// 余白は行頭の SGR より前に出すので、背景色付きの行でも余白は塗られない。
// 空行には余白を付けない。
func PrintLines(left int, lines []string) {
	pad := strings.Repeat(" ", max(0, left))
	var b strings.Builder
	for _, l := range lines {
		if l != "" {
			b.WriteString(pad)
			b.WriteString(l)
		}
		b.WriteString("\n")
	}
	backend.WriteString(b.String())
}
