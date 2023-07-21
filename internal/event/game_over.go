package event

import (
	"fmt"
	"jailbreak/pkg/system"
	"time"
)

func GameOver() {
	system.System("clear")

    fmt.Println(" #####      #     #     #  #######  #######  #     #  #######  ######  \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#     #    # #    ##   ##  #        #     #  #     #  #        #     # \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#         #   #   # # # #  #        #     #  #     #  #        #     # \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#  ####  #     #  #  #  #  #####    #     #  #     #  #####    ######  \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#     #  #######  #     #  #        #     #   #   #   #        #   #   \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#     #  #     #  #     #  #        #     #    # #    #        #    #  \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println(" #####   #     #  #     #  #######  #######     #     #######  #     # \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("")
    time.Sleep(500 * time.Millisecond)
    fmt.Println( "ざんねん！！わたしの　ぼうけんは　これで　おわってしまった！！\n" );
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\nPress [any key] to continue.")

	WaitToPressAnyKey()
}

func printGameOver() {
	system.System("clear")

    fmt.Println(" #####      #     #     #  #######  #######  #     #  #######  ######  \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#     #    # #    ##   ##  #        #     #  #     #  #        #     # \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#         #   #   # # # #  #        #     #  #     #  #        #     # \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#  ####  #     #  #  #  #  #####    #     #  #     #  #####    ######  \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#     #  #######  #     #  #        #     #   #   #   #        #   #   \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("#     #  #     #  #     #  #        #     #    # #    #        #    #  \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println(" #####   #     #  #     #  #######  #######     #     #######  #     # \n");
	time.Sleep(500 * time.Millisecond)
    fmt.Println("")
    time.Sleep(500 * time.Millisecond)
    fmt.Println( "ざんねん！！わたしの　ぼうけんは　これで　おわってしまった！！\n" );
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("\nPress [any key] to continue.")
}
