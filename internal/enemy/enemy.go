package enemy

import (
	"fmt"
	"jailbreak/pkg/system"
	"math/rand"
)

const typingTextFilePath = "assets/battle/keibi.txt"

type Enemy struct {
	Name  string
	FilePath string
	TextJapanese string
	TextRomaji string
}

func NewEnemy(name string, path string) *Enemy {
	line := 2 * (rand.Intn(56)) + 1 // ランダムで奇数行を選択
	textJapanese := system.LoadLineText(typingTextFilePath, line)
	textRomaji := system.LoadLineText(typingTextFilePath, line + 1)
	return &Enemy{Name: name, FilePath: path, TextJapanese: textJapanese, TextRomaji: textRomaji}
}

func (n *Enemy) Display() {
	system.PrintFile(n.FilePath)
}

func (n *Enemy) PrintTypingText() {
	fmt.Println(n.TextJapanese)
	fmt.Println(n.TextRomaji)
}