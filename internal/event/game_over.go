package event

import (
	"time"

	"jailbreak/internal/console"
)

func GameOver() {
	console.Clear()

	console.Printf(" #####      #     #     #  #######  #######  #     #  #######  ######  \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#     #    # #    ##   ##  #        #     #  #     #  #        #     # \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#         #   #   # # # #  #        #     #  #     #  #        #     # \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#  ####  #     #  #  #  #  #####    #     #  #     #  #####    ######  \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#     #  #######  #     #  #        #     #   #   #   #        #   #   \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#     #  #     #  #     #  #        #     #    # #    #        #    #  \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf(" #####   #     #  #     #  #######  #######     #     #######  #     # \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Println("")
	time.Sleep(500 * time.Millisecond)
	console.Printf("ざんねん！！わたしの　ぼうけんは　これで　おわってしまった！！\n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("\nPress [any key] to continue.")

	WaitToPressAnyKey()
}

func printGameOver() {
	console.Clear()

	console.Printf(" #####      #     #     #  #######  #######  #     #  #######  ######  \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#     #    # #    ##   ##  #        #     #  #     #  #        #     # \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#         #   #   # # # #  #        #     #  #     #  #        #     # \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#  ####  #     #  #  #  #  #####    #     #  #     #  #####    ######  \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#     #  #######  #     #  #        #     #   #   #   #        #   #   \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("#     #  #     #  #     #  #        #     #    # #    #        #    #  \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf(" #####   #     #  #     #  #######  #######     #     #######  #     # \n\n")
	time.Sleep(500 * time.Millisecond)
	console.Println("")
	time.Sleep(500 * time.Millisecond)
	console.Printf("ざんねん！！わたしの　ぼうけんは　これで　おわってしまった！！\n\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("\nPress [any key] to continue.")
}
