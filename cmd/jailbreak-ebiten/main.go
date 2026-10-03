// Ebitengine製フロントエンドのエントリポイント。
// デスクトップウィンドウと WASM(ブラウザ) の両方で動作する。
//
// 既存のゲームロジックは仮想端末(console.VTerm)に対して入出力を行い、
// 本フロントエンドは毎フレーム、仮想端末の文字グリッドを等幅フォントで描画し、
// キーボード入力を仮想端末の入力チャネルへ流すだけの構成とする。
package main

import (
	_ "embed"
	"image"
	"image/color"
	"image/draw"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"jailbreak/internal/console"
	"jailbreak/internal/game"
)

// 日本語等幅フォント JailbreakMono。
// HackGen Console Regular v2.10.0 (SIL Open Font License 1.1) を、ゲームが表示する
// 文字だけにサブセット化して改名したもの(OFL の Reserved Font Name 対応)。
// 生成手順は `make font`(tools/font/subset.py)、ライセンスは fonts/LICENSE_HackGen を参照。
//
//go:embed fonts/JailbreakMono-Regular.ttf
var fontTTF []byte

const (
	// 仮想端末のサイズ(推奨ターミナルサイズ 80x55 に合わせる)
	cols = 80
	rows = 55

	// フォントサイズと1セルの大きさ
	// HackGen Console は半角:全角 = 1:2 の等幅フォントなので
	// 半角セル幅 = フォントサイズ / 2 でグリッドが揃う
	fontSize = 16
	cellW    = fontSize / 2
	cellH    = fontSize

	screenW = cols * cellW // 640
	screenH = rows * cellH // 880
)

// ANSI 8色のカラーパレット
var ansiPalette = [8]color.RGBA{
	{0x00, 0x00, 0x00, 0xff}, // 0: 黒
	{0xcd, 0x31, 0x31, 0xff}, // 1: 赤
	{0x0d, 0xbc, 0x79, 0xff}, // 2: 緑
	{0xe5, 0xe5, 0x10, 0xff}, // 3: 黄
	{0x24, 0x72, 0xc8, 0xff}, // 4: 青
	{0xbc, 0x3f, 0xbc, 0xff}, // 5: マゼンタ
	{0x11, 0xa8, 0xcd, 0xff}, // 6: シアン
	{0xe5, 0xe5, 0xe5, 0xff}, // 7: 白
}

var (
	defaultFg = color.RGBA{0xe5, 0xe5, 0xe5, 0xff} // デフォルト文字色
	defaultBg = color.RGBA{0x0c, 0x0c, 0x0c, 0xff} // デフォルト背景色
)

// App はEbitengineのゲームインターフェース実装。
type App struct {
	vt    *console.VTerm
	glyph *glyphCache

	canvas    *image.RGBA   // 文字グリッドをCPUで描く先(内容が変わったときだけ再描画)
	offscreen *ebiten.Image // canvas を転送した画面用イメージ
	lastGen   int64
	cellBuf   []console.Cell
	done      chan struct{} // ゲームロジック終了通知
	inputBuf  []rune
}

// newFace は埋め込みフォントから描画用のフェイスを作る。
func newFace() font.Face {
	ft, err := opentype.Parse(fontTTF)
	if err != nil {
		log.Fatal(err)
	}
	face, err := opentype.NewFace(ft, &opentype.FaceOptions{Size: fontSize, DPI: 72, Hinting: font.HintingNone})
	if err != nil {
		log.Fatal(err)
	}
	return face
}

func NewApp() *App {
	vt := console.NewVTerm(cols, rows)
	console.SetBackend(vt)

	a := &App{
		vt:      vt,
		glyph:   newGlyphCache(newFace()),
		canvas:  image.NewRGBA(image.Rect(0, 0, screenW, screenH)),
		done:    make(chan struct{}),
		lastGen: -1,
	}

	// 既存のゲームロジックはブロッキングな逐次処理のまま、
	// 専用のgoroutineで走らせる
	go func() {
		game.Run()
		close(a.done)
	}()

	return a
}

func (a *App) Update() error {
	// ゲームロジックが終了(タイトルで「ゲーム終了」選択)したらアプリも終了
	select {
	case <-a.done:
		return ebiten.Termination
	default:
	}

	// 文字入力を仮想端末の入力チャネルへ流す
	a.inputBuf = ebiten.AppendInputChars(a.inputBuf[:0])
	for _, r := range a.inputBuf {
		a.vt.PushRune(r)
	}

	// Enterキー(ゲーム内ではコード13で判定される)
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
		a.vt.PushRune(13)
	}

	// 矢印キーも w/a/s/d として扱う(操作補助)
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		a.vt.PushRune('w')
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		a.vt.PushRune('s')
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		a.vt.PushRune('a')
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		a.vt.PushRune('d')
	}

	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	if a.offscreen == nil {
		a.offscreen = ebiten.NewImage(screenW, screenH)
	}

	// 仮想端末の内容が変わったときだけグリッドを再描画する
	if gen := a.vt.Gen(); gen != a.lastGen {
		a.renderGrid()
		a.lastGen = gen
	}

	screen.DrawImage(a.offscreen, nil)
}

func (a *App) renderGrid() {
	a.cellBuf = a.vt.Snapshot(a.cellBuf)
	drawGrid(a.canvas, a.cellBuf, a.glyph)
	a.offscreen.WritePixels(a.canvas.Pix)
}

// drawGrid は仮想端末の文字グリッド cells を dst に描く。
// Ebitengine の text パッケージは go-text/typesetting(HarfBuzz 移植等)を含み
// WASM が約3MB大きくなるため、golang.org/x/image/font でCPU描画して転送している。
func drawGrid(dst *image.RGBA, cells []console.Cell, g *glyphCache) {
	draw.Draw(dst, dst.Bounds(), image.NewUniform(defaultBg), image.Point{}, draw.Src)

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			c := cells[y*cols+x]
			if c.R == 0 { // 全角文字の2マス目
				continue
			}

			fg, bg := cellColors(c)

			// 背景色(デフォルト以外のときだけ塗る)
			if bg != defaultBg {
				w := cellW
				if x+1 < cols && cells[y*cols+x+1].R == 0 {
					w = cellW * 2 // 全角文字は2セル分
				}
				r := image.Rect(x*cellW, y*cellH, x*cellW+w, (y+1)*cellH)
				draw.Draw(dst, r, image.NewUniform(bg), image.Point{}, draw.Src)
			}

			// 空白は描画しない。HackGen Console は全角スペース(U+3000)を
			// 点線の四角で可視化するため、AA中の全角スペースも飛ばす
			if c.R == ' ' || c.R == '\u3000' {
				continue
			}

			// 文字を描画(太字はわずかにずらして重ね描きで擬似的に表現)
			g.draw(dst, x*cellW, y*cellH, c.R, fg, false)
			if c.Bold {
				g.draw(dst, x*cellW, y*cellH, c.R, fg, true)
			}
		}
	}
}

// cellColors はセルの属性から実際の前景色・背景色を決める。
func cellColors(c console.Cell) (fg, bg color.RGBA) {
	fg = defaultFg
	if c.Fg >= 0 && c.Fg < 8 {
		fg = ansiPalette[c.Fg]
	}
	bg = defaultBg
	if c.Bg >= 0 && c.Bg < 8 {
		bg = ansiPalette[c.Bg]
	}
	if c.Reverse {
		fg, bg = bg, fg
	}
	return fg, bg
}

func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenW, screenH
}

func main() {
	ebiten.SetWindowTitle("jailbreak - ダンジョン脱出ゲーム")
	// 縦880pxはノートPCだと収まらないことがあるので少し縮めて表示
	ebiten.SetWindowSize(screenW*7/8, screenH*7/8)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(NewApp()); err != nil {
		log.Fatal(err)
	}
}
