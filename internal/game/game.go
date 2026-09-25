// Package game はゲーム全体のメインループを提供する。
// ターミナル版(cmd/jailbreak)とEbitengine版(cmd/jailbreak-ebiten)で共有される。
package game

import (
	"time"

	"jailbreak/internal/console"
	"jailbreak/internal/dungeon"
	"jailbreak/internal/event"
	"jailbreak/internal/player"
	"jailbreak/internal/title"
)

// Run はゲームのメインループを実行する。
// タイトルで「ゲーム終了」が選択されると戻る。
// 呼び出し前に console.SetBackend で入出力バックエンドを設定しておくこと。
func Run() {
	for {
		// タイトル表示
		t := title.NewTitle()
		t.Select()

		switch t.Selected {
		case title.Rule:
			title.PrintRule()
			break
		case title.Start:
			break
		case title.End:
			console.Println("ゲームを終了します...")
			return
		}

		console.Println("Game Start!")

		p := player.NewPlayer()
		d := dungeon.Create(*p)

		// ゲームループ：勝利条件を満たすまで続く
	GameLoop:
		for {
			d.Display() // 部屋の様子を表示

			switch d.CheckEvent() {
			case event.GameClearEvent:
				event.GameClear()
				break GameLoop
			case event.GameOverEvent:
				event.GameOver()
				break GameLoop
			}

			d.WaitAction() // ユーザーから行動を入力してもらい、移動または方向転換を行う
			console.Printf("\n")

			time.Sleep(400 * time.Millisecond) // 0.4秒停止
		}
	}
}
