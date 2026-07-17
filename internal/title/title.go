package title

import (
	"jailbreak/internal/console"
)

type Selection int

const (
	Rule Selection = iota
	Start
	End
)

type Title struct {
	Selected Selection
}

func NewTitle() *Title {
	t := new(Title)
	t.Selected = 0
	return t
}

func (t *Title) Select() {
	t.print()

	for {
		r := console.ReadKey()
		switch r {
		case 119: // w
			t.prev()
		case 115: // s
			t.next()
		case 13: // Enter
			return
		}
		t.print()
	}
}

func (t *Title) next() {
	t.Selected++
	if t.Selected > 2 {
		t.Selected = 2
	}
}

func (t *Title) prev() {
	t.Selected--
	if t.Selected < 0 {
		t.Selected = 0
	}
}

func (t *Title) print() {
	console.Clear()

	console.PrintFile("assets/title/logo")

	switch t.Selected {
	case Rule:
		console.Printf("                       ")
		console.Printf("\x1b[7m")
		console.Printf("1) ルール説明\n")
		console.Printf("\x1b[0m")
		console.Printf("                       ")
		console.Printf("2) ゲームスタート\n")
		console.Printf("                       ")
		console.Printf("3) ゲーム終了\n")
	case Start:
		console.Printf("                       ")
		console.Printf("1) ルール説明\n")
		console.Printf("                       ")
		console.Printf("\x1b[7m")
		console.Printf("2) ゲームスタート\n")
		console.Printf("\x1b[0m")
		console.Printf("                       ")
		console.Printf("3) ゲーム終了\n")
	case End:
		console.Printf("                       ")
		console.Printf("1) ルール説明\n")
		console.Printf("                       ")
		console.Printf("2) ゲームスタート\n")
		console.Printf("                       ")
		console.Printf("\x1b[7m")
		console.Printf("3) ゲーム終了\n")
		console.Printf("\x1b[0m")
	}
}
