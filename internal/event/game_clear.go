package event

import (
	"time"

	"jailbreak/internal/console"
)

func GameClear() {
	console.Clear()
	console.Printf("\nゲームクリア！！\n")
	time.Sleep(1 * time.Second)
	console.Printf("\n")
	console.Printf("Thank you for playing.\n")
	time.Sleep(1 * time.Second)
	console.Printf("\n")
	console.Printf("\n")
	console.Printf("Produced by raku2wei")
	console.Printf("\n")

	console.Printf("\nPress [any key] to continue.")

	WaitToPressAnyKey()
}
