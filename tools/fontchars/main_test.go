package main

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/image/font/sfnt"
)

const (
	fontPath  = "../../cmd/jailbreak-ebiten/fonts/JailbreakMono-Regular.ttf"
	charsPath = "../../cmd/jailbreak-ebiten/fonts/chars.txt"
)

// 画面に出す文字を追加したのにフォントを再生成し忘れると、
// ブラウザ版・デスクトップ版でその文字が豆腐(□)になる。それをテストで検出する。
func TestEmbeddedFontCoversGameText(t *testing.T) {
	chars, err := collect("../..")
	if err != nil {
		t.Fatal(err)
	}

	b, err := os.ReadFile(fontPath)
	if err != nil {
		t.Fatal(err)
	}
	f, err := sfnt.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	var buf sfnt.Buffer
	var missing []rune
	for _, r := range chars {
		if r == ' ' || r == '　' { // 空白は描画しない
			continue
		}
		if gi, err := f.GlyphIndex(&buf, r); err != nil || gi == 0 {
			missing = append(missing, r)
		}
	}
	if len(missing) > 0 {
		t.Errorf("埋め込みフォントに無い文字があります。`make font` で再生成してください: %q", string(missing))
	}

	// chars.txt も最新であること(差分レビューのためにコミットしている)
	want, err := os.ReadFile(charsPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimRight(string(want), "\n") != chars {
		t.Errorf("%s が古いです。`make font` で再生成してください", charsPath)
	}
}
