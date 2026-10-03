package enemy

import (
	"math/rand"

	"jailbreak/assets"
	"jailbreak/internal/console"
)

const typingTextFilePath = "assets/battle/keibi.txt"

type Enemy struct {
	Name         string
	FilePath     string
	TextJapanese string
	TextRomaji   string
}

func NewEnemy(name string, path string) *Enemy {
	line := 2*(rand.Intn(56)) + 1 // ランダムで奇数行を選択
	textJapanese := assets.LoadLineText(typingTextFilePath, line)
	textRomaji := assets.LoadLineText(typingTextFilePath, line+1)
	return &Enemy{Name: name, FilePath: path, TextJapanese: textJapanese, TextRomaji: textRomaji}
}

func (n *Enemy) Display() {
	console.PrintFile(n.FilePath)
}

func (n *Enemy) PrintTypingText() {
	console.Println(n.TextJapanese)
	console.Println(n.TextRomaji)
}
