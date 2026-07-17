package console

import (
	"sync"

	"github.com/mattn/go-runewidth"
)

// ColorDefault はデフォルト色(前景・背景)を表す。
// それ以外は ANSI SGR の 30-37 / 40-47 に対応する 0〜7 のカラーインデックス。
const ColorDefault = -1

// Cell は仮想スクリーン上の1マス。
// 全角文字はセル2つ分を占有し、2マス目には R=0 (継続マーク)が入る。
type Cell struct {
	R       rune
	Fg      int // ColorDefault または 0〜7
	Bg      int // ColorDefault または 0〜7
	Bold    bool
	Reverse bool
}

// VTerm は80x55等の文字グリッドを持つ簡易仮想端末。
// このゲームが使用する範囲のANSIエスケープ(SGR: 0,1,7,30-37,39,40-47,49)をパースして
// セル属性に反映する。Ebitengine 側は Snapshot() の結果を毎フレーム描画するだけでよい。
type VTerm struct {
	mu         sync.Mutex
	cols, rows int
	cells      []Cell
	cx, cy     int

	// 現在の描画属性
	fg, bg        int
	bold, reverse bool

	// エスケープシーケンスのパース状態
	inEsc bool
	esc   []rune

	gen   int64     // 画面内容の世代番号(再描画判定用)
	input chan rune // キー入力チャネル
}

// NewVTerm は cols x rows の仮想端末を生成する。
func NewVTerm(cols, rows int) *VTerm {
	v := &VTerm{
		cols:  cols,
		rows:  rows,
		cells: make([]Cell, cols*rows),
		fg:    ColorDefault,
		bg:    ColorDefault,
		input: make(chan rune, 64),
	}
	v.resetCells()
	return v
}

// Size は列数・行数を返す。
func (v *VTerm) Size() (cols, rows int) {
	return v.cols, v.rows
}

// Gen は画面内容の世代番号を返す。変化していなければ再描画不要と判断できる。
func (v *VTerm) Gen() int64 {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.gen
}

// Snapshot は現在のセル内容のコピーを buf に書き込んで返す。
// buf が nil または長さ不足の場合は新規確保する。
func (v *VTerm) Snapshot(buf []Cell) []Cell {
	v.mu.Lock()
	defer v.mu.Unlock()
	if len(buf) < len(v.cells) {
		buf = make([]Cell, len(v.cells))
	}
	copy(buf, v.cells)
	return buf
}

// PushRune はキー入力を入力チャネルへ送る(Ebitengineの Update から呼ぶ)。
// バッファが一杯の場合は取りこぼす(ゲームの性質上問題ない)。
func (v *VTerm) PushRune(r rune) {
	select {
	case v.input <- r:
	default:
	}
}

// ReadKey はキー入力を1文字待つ(ゲームロジック側のgoroutineから呼ばれる)。
func (v *VTerm) ReadKey() rune {
	return <-v.input
}

// TryReadKey は入力があれば取り出す(ノンブロッキング)。
func (v *VTerm) TryReadKey() (rune, bool) {
	select {
	case r := <-v.input:
		return r, true
	default:
		return 0, false
	}
}

// Clear は画面全体をクリアしてカーソルを左上に戻す。
func (v *VTerm) Clear() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.resetCells()
	v.cx, v.cy = 0, 0
	v.gen++
}

// WriteString はANSIエスケープを含み得る文字列を仮想スクリーンへ反映する。
func (v *VTerm) WriteString(s string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	for _, r := range s {
		v.writeRune(r)
	}
	v.gen++
}

func (v *VTerm) resetCells() {
	for i := range v.cells {
		v.cells[i] = Cell{R: ' ', Fg: ColorDefault, Bg: ColorDefault}
	}
}

func (v *VTerm) writeRune(r rune) {
	// エスケープシーケンスのパース中
	if v.inEsc {
		if len(v.esc) == 0 {
			if r == '[' { // CSIシーケンス開始
				v.esc = append(v.esc, r)
			} else {
				// CSI以外のエスケープは未対応なので読み捨てる
				v.inEsc = false
			}
			return
		}
		if r >= 0x40 && r <= 0x7e { // 終端文字
			if r == 'm' {
				v.applySGR(string(v.esc[1:]))
			}
			// SGR以外のCSI(カーソル移動等)はこのゲームでは未使用なので無視
			v.inEsc = false
			v.esc = v.esc[:0]
			return
		}
		v.esc = append(v.esc, r)
		return
	}

	switch r {
	case 0x1b: // ESC
		v.inEsc = true
		v.esc = v.esc[:0]
	case '\n':
		v.newline()
	case '\r':
		v.cx = 0
	case '\t':
		// タブは8桁区切りに丸める
		v.cx = (v.cx/8 + 1) * 8
		if v.cx >= v.cols {
			v.newline()
		}
	default:
		v.putRune(r)
	}
}

func (v *VTerm) putRune(r rune) {
	w := runewidth.RuneWidth(r)
	if w <= 0 {
		return // 制御文字・結合文字は無視
	}
	if v.cx+w > v.cols {
		v.newline()
	}
	c := Cell{R: r, Fg: v.fg, Bg: v.bg, Bold: v.bold, Reverse: v.reverse}
	v.cells[v.cy*v.cols+v.cx] = c
	if w == 2 && v.cx+1 < v.cols {
		// 全角文字の2マス目は継続マーク
		v.cells[v.cy*v.cols+v.cx+1] = Cell{R: 0, Fg: v.fg, Bg: v.bg, Bold: v.bold, Reverse: v.reverse}
	}
	v.cx += w
}

func (v *VTerm) newline() {
	v.cx = 0
	v.cy++
	if v.cy >= v.rows {
		// 1行スクロール
		copy(v.cells, v.cells[v.cols:])
		for i := (v.rows - 1) * v.cols; i < len(v.cells); i++ {
			v.cells[i] = Cell{R: ' ', Fg: ColorDefault, Bg: ColorDefault}
		}
		v.cy = v.rows - 1
	}
}

// applySGR は "1;36" のようなSGRパラメータ文字列を現在属性に反映する。
func (v *VTerm) applySGR(params string) {
	if params == "" {
		params = "0"
	}
	n := 0
	apply := func(code int) {
		switch {
		case code == 0:
			v.fg, v.bg = ColorDefault, ColorDefault
			v.bold, v.reverse = false, false
		case code == 1:
			v.bold = true
		case code == 7:
			v.reverse = true
		case code == 22:
			v.bold = false
		case code == 27:
			v.reverse = false
		case code >= 30 && code <= 37:
			v.fg = code - 30
		case code == 39:
			v.fg = ColorDefault
		case code >= 40 && code <= 47:
			v.bg = code - 40
		case code == 49:
			v.bg = ColorDefault
		}
	}
	for _, ch := range params {
		if ch >= '0' && ch <= '9' {
			n = n*10 + int(ch-'0')
		} else if ch == ';' {
			apply(n)
			n = 0
		}
	}
	apply(n)
}
