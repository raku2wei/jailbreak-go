package event

import (
	"jailbreak/internal/console"
)

type Event int

const (
	NoEvent Event = iota
	GameClearEvent
	GameOverEvent
)

func WaitToPressAnyKey() {
	// 任意のキー入力を待つ
	console.ReadKey()
}
