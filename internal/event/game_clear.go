package event

import (
	"fmt"
	"jailbreak/pkg/system"
	"log"
	"time"

	"github.com/mattn/go-tty"
)

func GameClear() {
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
}
