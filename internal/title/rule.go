package title

import (
	"jailbreak/internal/console"
)

func PrintRule() {
	print()

	// 任意のキー入力を待つ
	console.ReadKey()
}

func print() {
	console.Clear()
	console.PrintFile("assets/title/rule")
}
