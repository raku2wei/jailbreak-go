package main

import (
	"image"
	"image/color"
	"image/draw"

	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// boldOffset は擬似太字の重ね描きのずらし量(1/64 px 単位)。
// 以前の text/v2 描画で 0.7px を指定したときの実効値(1/8px 単位に丸めた 0.625px)に合わせている。
const boldOffset = fixed.Int26_6(40)

// glyphCache は文字ごとのラスタライズ結果(アルファマスク)をキャッシュする。
// Ebitengine の text パッケージ(go-text/typesetting を含む)を使わずに
// golang.org/x/image/font で直接描くことで WASM のサイズを抑えている。
type glyphCache struct {
	face   font.Face
	ascent fixed.Int26_6
	m      map[glyphKey]*glyphMask
}

type glyphKey struct {
	r    rune
	bold bool
}

type glyphMask struct {
	mask *image.Alpha
	off  image.Point // セル左上からのマスク左上の位置
}

func newGlyphCache(face font.Face) *glyphCache {
	// ベースラインは以前の text/v2 描画と同じく整数ピクセルに切り捨てる(文字がにじまない)
	return &glyphCache{face: face, ascent: fixed.I(face.Metrics().Ascent.Floor()), m: map[glyphKey]*glyphMask{}}
}

// draw はセル左上 (x, y) を基準に文字 r を色 c で dst に描く。
func (g *glyphCache) draw(dst *image.RGBA, x, y int, r rune, c color.RGBA, bold bool) {
	k := glyphKey{r, bold}
	gm, ok := g.m[k]
	if !ok {
		gm = g.rasterize(r, bold)
		g.m[k] = gm
	}
	if gm == nil {
		return
	}
	b := gm.mask.Bounds()
	p := image.Pt(x, y).Add(gm.off)
	draw.DrawMask(dst, image.Rectangle{Min: p, Max: p.Add(b.Size())}, image.NewUniform(c), image.Point{}, gm.mask, b.Min, draw.Over)
}

func (g *glyphCache) rasterize(r rune, bold bool) *glyphMask {
	dot := fixed.Point26_6{Y: g.ascent}
	if bold {
		dot.X = boldOffset
	}
	// フォントに無い文字は ok=false でも .notdef(豆腐 □)のマスクが返るので、
	// それをそのまま描いて欠けに気付けるようにする
	dr, mask, mp, _, _ := g.face.Glyph(dot, r)
	if mask == nil || dr.Empty() {
		return nil
	}
	// face.Glyph が返すマスクは内部バッファを使い回すため複製しておく
	a := image.NewAlpha(image.Rect(0, 0, dr.Dx(), dr.Dy()))
	draw.Draw(a, a.Bounds(), mask, mp, draw.Src)
	return &glyphMask{mask: a, off: dr.Min}
}
