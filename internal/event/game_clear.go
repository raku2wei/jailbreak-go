package event

import (
	"fmt"
	"jailbreak/pkg/system"
	"time"
)

func GameClear() {
	system.System("clear")
	fmt.Printf("\nゲームクリア！！\n")
	time.Sleep(1 * time.Second)
	fmt.Printf("\n")
	fmt.Printf("Thank you for playing.\n")
	time.Sleep(1 * time.Second)
	fmt.Printf("\n")
	fmt.Printf("\n")
	fmt.Printf("Produced by raku2wei")
	fmt.Printf("\n")

	fmt.Printf("\nPress [any key] to continue.")

	WaitToPressAnyKey()
}
