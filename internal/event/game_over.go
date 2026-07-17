package event

import (
	"time"

	"jailbreak/internal/console"
)

func GameOver() {
	console.Clear()

	console.Println(" #####      #     #     #  #######  #######  #     #  #######  ######  \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#     #    # #    ##   ##  #        #     #  #     #  #        #     # \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#         #   #   # # # #  #        #     #  #     #  #        #     # \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#  ####  #     #  #  #  #  #####    #     #  #     #  #####    ######  \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#     #  #######  #     #  #        #     #   #   #   #        #   #   \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#     #  #     #  #     #  #        #     #    # #    #        #    #  \n")
	time.Sleep(500 * time.Millisecond)
	console.Println(" #####   #     #  #     #  #######  #######     #     #######  #     # \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("")
	time.Sleep(500 * time.Millisecond)
	console.Println("ざんねん！！わたしの　ぼうけんは　これで　おわってしまった！！\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("\nPress [any key] to continue.")

	WaitToPressAnyKey()
}

func printGameOver() {
	console.Clear()

	console.Println(" #####      #     #     #  #######  #######  #     #  #######  ######  \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#     #    # #    ##   ##  #        #     #  #     #  #        #     # \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#         #   #   # # # #  #        #     #  #     #  #        #     # \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#  ####  #     #  #  #  #  #####    #     #  #     #  #####    ######  \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#     #  #######  #     #  #        #     #   #   #   #        #   #   \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("#     #  #     #  #     #  #        #     #    # #    #        #    #  \n")
	time.Sleep(500 * time.Millisecond)
	console.Println(" #####   #     #  #     #  #######  #######     #     #######  #     # \n")
	time.Sleep(500 * time.Millisecond)
	console.Println("")
	time.Sleep(500 * time.Millisecond)
	console.Println("ざんねん！！わたしの　ぼうけんは　これで　おわってしまった！！\n")
	time.Sleep(500 * time.Millisecond)
	console.Printf("\nPress [any key] to continue.")
}
