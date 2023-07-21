package event

import (
	"log"

	"github.com/mattn/go-tty"
)

type Event int

const (
	NoEvent Event = iota
	GameClearEvent
	GameOverEvent
)

func WaitToPressAnyKey() {
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