package enemy

import (
	"fmt"
	"jailbreak/internal/event"
	"jailbreak/pkg/system"
	"log"
	"time"

	"github.com/mattn/go-tty"
)

const battleTimeLimit = 15 * time.Second

// 通常戦闘処理
func (e *Enemy) Battle() bool {

	time.Sleep(500 * time.Millisecond)

	system.System("clear")

	fmt.Println(e.Name + "があらわれた！")
	fmt.Println("")
	fmt.Printf("Press [any key] to Start")

	event.WaitToPressAnyKey()

	var isMiss bool = false      // ミス判定用 

	var startTime = time.Now()

	var answer []rune = []rune(e.TextRomaji)
	var playerInput [128]rune
	var m rune = 0; 
	var n int = 0;

	inputChan := make(chan rune)

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer tty.Close()

	// 画面描写をブロックしないように、入力の受付はgoroutineで実施
	go func() {
		for {
			r, err := tty.ReadRune()
			if err != nil {
				log.Fatal(err)
			}

			inputChan <- r
		}
	}()

	for n < len(answer) {
		select {
		case r := <-inputChan:
			if isMiss {
				isMiss = false  // ミス判定解除
			}
			if answer[n] == r {			// 問題の文字と入力した文字が一致したら
				playerInput[n] = r		// 入力した文字を回答用変数に代入して
				playerInput[n+1] = '\x00'; // 終端文字を追加
				n++;
			} else {					// 一致しなかったらミス判定
				isMiss = true;
				m = r;
			}

			if !refresh(e, isMiss, string(playerInput[:n]), m, startTime) {
				return false
			}
		default:
			if !refresh(e, isMiss, string(playerInput[:n]), m, startTime) {
				return false
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	fmt.Println("\nYou Win!!")
	time.Sleep(1000 * time.Millisecond)

	system.System("clear")

	return true
}

func refresh(e *Enemy, isMiss bool, playerInput string, m rune, startTime time.Time ) bool {
	system.System("clear")        // 画面クリア
		
	// 敵のAAを表示
	e.Display()

	fmt.Printf("\x1b[1m")
	fmt.Printf("「%s」\n", e.TextJapanese)    // 問題文(日本語)表示
	fmt.Printf("%s\n", e.TextRomaji)     // 問題文(ローマ字)表示
	fmt.Printf("\x1b[0m")
  
	fmt.Printf("\x1b[32m") // 文字色を緑に変更
	fmt.Printf("%s", playerInput)  // 入力された文字列を表示
	fmt.Printf("\x1b[39m")     // 文字色を戻す


	if isMiss {
		fmt.Printf("\x1b[41m")         // 文字背景色を赤に変更
		fmt.Printf("%s\n", string(m))          // 入力した文字表示
		fmt.Printf("タイプミス！\n")   // ミスメッセージ表示
		fmt.Printf("\x1b[49m")         // 文字背景色を戻す
	} else {
		fmt.Printf("\n\n")
	}
	// r = 0;           // rクリア
	timeLimit := battleTimeLimit.Seconds() - time.Now().Sub(startTime).Seconds()      // 制限時間(残り時間)計算
	if timeLimit <= 0 {
		timeLimit = 0
	}
	fmt.Printf("残り時間 : %.2f 秒\n", timeLimit)       // 制限時間表示
	if timeLimit <= 0 {                           // 時間切れになったらゲームオーバー
		fmt.Printf("TIME OVER!\n");
		return false
	}

	return true
}
