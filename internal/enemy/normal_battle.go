package enemy

import (
	"time"

	"jailbreak/internal/console"
	"jailbreak/internal/event"
)

const battleTimeLimit = 15 * time.Second

// 通常戦闘処理
func (e *Enemy) Battle() bool {

	time.Sleep(500 * time.Millisecond)

	console.Clear()

	console.Println(e.Name + "があらわれた！")
	console.Println("")
	console.Printf("Press [any key] to Start")

	event.WaitToPressAnyKey()

	var isMiss bool = false // ミス判定用

	var startTime = time.Now()

	var answer []rune = []rune(e.TextRomaji)
	var playerInput [128]rune
	var m rune = 0
	var n int = 0

	for n < len(answer) {
		// 画面描写をブロックしないように、入力はノンブロッキングで受け付ける
		if r, ok := console.TryReadKey(); ok {
			if isMiss {
				isMiss = false // ミス判定解除
			}
			if answer[n] == r { // 問題の文字と入力した文字が一致したら
				playerInput[n] = r        // 入力した文字を回答用変数に代入して
				playerInput[n+1] = '\x00' // 終端文字を追加
				n++
			} else { // 一致しなかったらミス判定
				isMiss = true
				m = r
			}

			if !refresh(e, isMiss, string(playerInput[:n]), m, startTime) {
				return false
			}
		} else {
			if !refresh(e, isMiss, string(playerInput[:n]), m, startTime) {
				return false
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	console.Println("\nYou Win!!")
	time.Sleep(1000 * time.Millisecond)

	console.Clear()

	return true
}

func refresh(e *Enemy, isMiss bool, playerInput string, m rune, startTime time.Time) bool {
	console.Clear() // 画面クリア

	// 敵のAAを表示
	e.Display()

	console.Printf("\x1b[1m")
	console.Printf("「%s」\n", e.TextJapanese) // 問題文(日本語)表示
	console.Printf("%s\n", e.TextRomaji)     // 問題文(ローマ字)表示
	console.Printf("\x1b[0m")

	console.Printf("\x1b[32m")        // 文字色を緑に変更
	console.Printf("%s", playerInput) // 入力された文字列を表示
	console.Printf("\x1b[39m")        // 文字色を戻す

	if isMiss {
		console.Printf("\x1b[41m")        // 文字背景色を赤に変更
		console.Printf("%s\n", string(m)) // 入力した文字表示
		console.Printf("タイプミス！\n")        // ミスメッセージ表示
		console.Printf("\x1b[49m")        // 文字背景色を戻す
	} else {
		console.Printf("\n\n")
	}
	// r = 0;           // rクリア
	timeLimit := battleTimeLimit.Seconds() - time.Now().Sub(startTime).Seconds() // 制限時間(残り時間)計算
	if timeLimit <= 0 {
		timeLimit = 0
	}
	console.Printf("残り時間 : %.2f 秒\n", timeLimit) // 制限時間表示
	if timeLimit <= 0 {                          // 時間切れになったらゲームオーバー
		console.Printf("TIME OVER!\n")
		return false
	}

	return true
}
