package title

import (
	"fmt"
	"jailbreak/pkg/system"
	"log"

	"github.com/mattn/go-tty"
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

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer tty.Close()

	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
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
	system.System("clear")

	system.PrintFile("assets/title/logo")

	switch t.Selected {
	case Rule:
		fmt.Printf("                       ")
		fmt.Printf("\x1b[7m")
		fmt.Printf("1) ルール説明\n")
		fmt.Printf("\x1b[0m")
		fmt.Printf("                       ")
		fmt.Printf("2) ゲームスタート\n")
		fmt.Printf("                       ")
		fmt.Printf("3) ゲーム終了\n")
	case Start:
		fmt.Printf("                       ")
		fmt.Printf("1) ルール説明\n")
		fmt.Printf("                       ")
		fmt.Printf("\x1b[7m")
		fmt.Printf("2) ゲームスタート\n")
		fmt.Printf("\x1b[0m")
		fmt.Printf("                       ")
		fmt.Printf("3) ゲーム終了\n")
	case End:
		fmt.Printf("                       ")
		fmt.Printf("1) ルール説明\n")
		fmt.Printf("                       ")
		fmt.Printf("2) ゲームスタート\n")
		fmt.Printf("                       ")
		fmt.Printf("\x1b[7m")
		fmt.Printf("3) ゲーム終了\n")
		fmt.Printf("\x1b[0m")
	}
}
