package enemy

import (
	"fmt"
	"strings"
	"time"

	"jailbreak/assets"
	"jailbreak/internal/console"
	"jailbreak/internal/event"
)

const battleTimeLimit = 15 * time.Second

// 戦闘画面に出す定型の文言
const (
	startMessage    = "Press [any key] to Start"
	missMessage     = "タイプミス！"
	timeLeftFormat  = "残り時間 : %.2f 秒"
	timeOverMessage = "TIME OVER!"
	winMessage      = "You Win!!"
)

const (
	// battleGapRows は敵のAAと問題文の間の空行数(従来のレイアウトと同じ)
	battleGapRows = 3
	// battleTextRows は問題文から勝敗表示までの行数
	// (問題文・ローマ字・入力・ミス表示・残り時間・空行・勝敗)
	battleTextRows = 7
)

// battleLayout は戦闘画面を画面(80x55)の中央に描くための余白。
// 入力の途中で位置が揺れないよう、戦闘開始時に戦闘中に変わらない内容から1回だけ計算する。
// 敵のAAと問題文以下はそれぞれの幅で左右中央に置き、全体の高さで上下中央に置く。
type battleLayout struct {
	top      int      // 画面上端から敵のAAまでの行数
	aa       []string // 敵のAA(ファイル前後の空行を除いたもの)
	aaLeft   int      // 敵のAAの左余白
	textLeft int      // 問題文以下の左余白(入力欄がローマ字の真下に来るよう共通)
}

func newBattleLayout(e *Enemy) battleLayout {
	aa := console.TrimBlankLines(strings.Split(string(assets.MustRead(e.FilePath)), "\n"))
	// 入力欄(入力済み文字 + ミスした1文字)はローマ字より長くならないので幅の計算に含めない
	textW := console.BlockWidth([]string{
		"「" + e.TextJapanese + "」",
		e.TextRomaji,
		missMessage,
		fmt.Sprintf(timeLeftFormat, battleTimeLimit.Seconds()),
		timeOverMessage,
		winMessage,
	})
	return battleLayout{
		top:      console.CenterTop(len(aa) + battleGapRows + battleTextRows),
		aa:       aa,
		aaLeft:   console.CenterLeft(console.BlockWidth(aa)),
		textLeft: console.CenterLeft(textW),
	}
}

// 通常戦闘処理
func (e *Enemy) Battle() bool {

	time.Sleep(500 * time.Millisecond)

	console.Clear()

	intro := []string{e.Name + "があらわれた！", "", startMessage}
	console.Printf("%s", strings.Repeat("\n", console.CenterTop(len(intro))))
	console.PrintLines(console.CenterLeft(console.BlockWidth(intro)), intro)

	event.WaitToPressAnyKey()

	layout := newBattleLayout(e)

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

			if !refresh(e, layout, isMiss, string(playerInput[:n]), m, startTime) {
				return false
			}
		} else {
			if !refresh(e, layout, isMiss, string(playerInput[:n]), m, startTime) {
				return false
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	console.PrintLines(layout.textLeft, []string{"", winMessage})
	time.Sleep(1000 * time.Millisecond)

	console.Clear()

	return true
}

func refresh(e *Enemy, l battleLayout, isMiss bool, playerInput string, m rune, startTime time.Time) bool {
	console.Clear() // 画面クリア

	// 上の余白
	console.Printf("%s", strings.Repeat("\n", l.top))

	// 敵のAAを表示
	console.PrintLines(l.aaLeft, l.aa)
	console.Printf("%s", strings.Repeat("\n", battleGapRows))

	lines := []string{
		"\x1b[1m「" + e.TextJapanese + "」\x1b[0m", // 問題文(日本語)
		"\x1b[1m" + e.TextRomaji + "\x1b[0m",     // 問題文(ローマ字)
	}

	input := "\x1b[32m" + playerInput + "\x1b[39m" // 入力された文字列を緑で表示
	if isMiss {
		// 入力した文字とミスメッセージを赤背景で表示
		lines = append(lines,
			input+"\x1b[41m"+string(m)+"\x1b[49m",
			"\x1b[41m"+missMessage+"\x1b[49m",
		)
	} else {
		lines = append(lines, input, "")
	}

	timeLimit := battleTimeLimit.Seconds() - time.Now().Sub(startTime).Seconds() // 制限時間(残り時間)計算
	if timeLimit <= 0 {
		timeLimit = 0
	}
	lines = append(lines, fmt.Sprintf(timeLeftFormat, timeLimit)) // 制限時間表示
	if timeLimit <= 0 {                                           // 時間切れになったらゲームオーバー
		lines = append(lines, timeOverMessage)
	}

	console.PrintLines(l.textLeft, lines)

	return timeLimit > 0
}
