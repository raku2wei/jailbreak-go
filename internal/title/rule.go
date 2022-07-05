package title

import (
	"jailbreak/pkg/system"
	"log"

	"github.com/mattn/go-tty"
)

func PrintRule() {
	print()

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer tty.Close()

	for {
		_, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		return
	}
}

func print() {
	system.System("clear")
	system.PrintFile("assets/title/rule")
}
