package dungeon

import (
	"fmt"
	"jailbreak/internal/enemy"
	"jailbreak/internal/event"
	"jailbreak/internal/player"
	"jailbreak/internal/room"
	"jailbreak/pkg/system"
	"log"
	"time"

	"github.com/mattn/go-tty"
)

type Dungeon struct {
	// ダンジョンの部屋は4x6のグリッドなので2次配列
	// (1,1)からスタートとするので配列は5x7
	rooms    [7][5]room.Room
	player   player.Player // ダンジョン内にいるプレイヤー
	moveRoom bool          // プレイヤーが部屋を移動したかどうか（イベントチェック用）
	encounterRate float32
}

// ダンジョン生成（初期化）
func Create(p player.Player) *Dungeon {
	d := new(Dungeon)

	// プレイヤーを配置
	// スタート位置は(1, 1)
	p.SetPosition(1, 1, room.North)
	d.player = p

	// ゴールを設置
	d.rooms[6][4].IsGoal = true

	// d.rooms[1][1].IsVisited= true
	// 各部屋のドアの情報
	// 1行目のドア情報
	d.rooms[1][1].HasDoor[room.East] = true
	d.rooms[2][1].HasDoor[room.North] = true
	d.rooms[2][1].HasDoor[room.West] = true
	d.rooms[3][1].HasDoor[room.East] = true
	d.rooms[3][1].HasDoor[room.West] = true
	d.rooms[4][1].HasDoor[room.North] = true
	d.rooms[4][1].HasDoor[room.East] = true
	d.rooms[4][1].HasDoor[room.West] = true
	d.rooms[5][1].HasDoor[room.North] = true
	d.rooms[5][1].HasDoor[room.East] = true
	d.rooms[5][1].HasDoor[room.West] = true
	d.rooms[6][1].HasDoor[room.North] = true
	d.rooms[6][1].HasDoor[room.West] = true

	// 2行目のドア情報
	d.rooms[1][2].HasDoor[room.East] = true
	d.rooms[2][2].HasDoor[room.North] = true
	d.rooms[2][2].HasDoor[room.South] = true
	d.rooms[2][2].HasDoor[room.West] = true
	d.rooms[3][2].HasDoor[room.North] = true
	d.rooms[3][2].HasDoor[room.East] = true
	d.rooms[4][2].HasDoor[room.South] = true
	d.rooms[4][2].HasDoor[room.West] = true
	d.rooms[5][2].HasDoor[room.North] = true
	d.rooms[5][2].HasDoor[room.East] = true
	d.rooms[5][2].HasDoor[room.South] = true
	d.rooms[6][2].HasDoor[room.North] = true
	d.rooms[6][2].HasDoor[room.South] = true
	d.rooms[6][2].HasDoor[room.West] = true

	// 3行目のドア情報
	d.rooms[1][3].HasDoor[room.North] = true
	d.rooms[1][3].HasDoor[room.East] = true
	d.rooms[2][3].HasDoor[room.South] = true
	d.rooms[2][3].HasDoor[room.West] = true
	d.rooms[3][3].HasDoor[room.South] = true
	d.rooms[3][3].HasDoor[room.West] = true
	d.rooms[4][3].HasDoor[room.North] = true
	d.rooms[5][3].HasDoor[room.North] = true
	d.rooms[5][3].HasDoor[room.East] = true
	d.rooms[5][3].HasDoor[room.South] = true
	d.rooms[6][3].HasDoor[room.North] = true
	d.rooms[6][3].HasDoor[room.South] = true
	d.rooms[6][3].HasDoor[room.West] = true

	// 4行目のドア情報
	d.rooms[1][4].HasDoor[room.East] = true
	d.rooms[1][4].HasDoor[room.South] = true
	d.rooms[2][4].HasDoor[room.East] = true
	d.rooms[2][4].HasDoor[room.South] = true
	d.rooms[3][4].HasDoor[room.West] = true
	d.rooms[3][4].HasDoor[room.East] = true
	d.rooms[3][4].HasDoor[room.South] = true
	d.rooms[4][4].HasDoor[room.West] = true
	d.rooms[4][4].HasDoor[room.East] = true
	d.rooms[4][4].HasDoor[room.South] = true
	d.rooms[5][4].HasDoor[room.West] = true
	d.rooms[6][4].HasDoor[room.South] = true

	//部屋にヒントを配置
	d.rooms[2][1].HasHint = true
	d.rooms[4][1].HasHint = true
	d.rooms[5][1].HasHint = true
	d.rooms[3][2].HasHint = true
	d.rooms[4][2].HasHint = true
	d.rooms[1][3].HasHint = true
	d.rooms[2][3].HasHint = true
	d.rooms[5][3].HasHint = true
	d.rooms[6][3].HasHint = true
	d.rooms[1][4].HasHint = true
	d.rooms[3][4].HasHint = true

	//罠(強制ワープ)部屋を設置
	d.rooms[6][2].IsWana = true

	// エンカウント率初期化
	d.encounterRate = DefaultEncounterRate

	return d
}

func (d *Dungeon) Display() {
	//プレーヤーの向き
	dir := d.player.Direction
	// プレイヤーが現在いる部屋を表示
	d.currentRoom().Display(dir)
}

func (d *Dungeon) CheckEvent() event.Event {

	fmt.Printf("\n")

	if !d.currentRoom().IsVisited {
		d.currentRoom().IsVisited = true
	}

	if d.moveRoom { // 部屋を移動したら判定
		if d.currentRoom().IsWana { // 罠があったら
			// 罠を起動
			d.wanaActivate()
		} else if d.currentRoom().IsGoal { // ゴールだったら
			return event.GameClearEvent
		} else if d.checkEncounter() {
			e := enemy.NewEnemy("警備員", "assets/enemy/keibi")
			if !e.Battle() {
				return event.GameOverEvent
			}
			// エンカウント率初期化
			d.encounterRate = DefaultEncounterRate
		} else {
			// エンカウント率増加
			d.encounterRate *= EncounterIncreaseRate
		}
		d.moveRoom = false
	}
	return event.NoEvent
}

func (d *Dungeon) WaitAction() {
	front := d.player.Direction // プレーヤー向き取得
	back := front + 2           // 後ろの方向を取得(バック用)
	if back > 3 {
		back -= 4
	}

	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}

	defer tty.Close()

	// 入力処理
	for {
		fmt.Printf("w：前に進む   s：後ろに進む   a：左を向く   d：右を向く   ")
		fmt.Printf("\x1b[1m")  // 文字を強調
		fmt.Printf("\x1b[36m") // 文字を水色に
		if front == room.North {
			fmt.Printf("↑ 北")
		} else if front == room.East {
			fmt.Printf("← 北")
		} else if front == room.South {
			fmt.Printf("↓ 北")
		} else if front == room.West {
			fmt.Printf("→ 北")
		} else {
			fmt.Printf("向きがおかしいよ")
		}
		fmt.Printf("\x1b[39m") // 文字色を戻す
		fmt.Printf("\x1b[0m")  // 強調を解除
		fmt.Printf("\n")

		d.PrintMap() // マップ表示

		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}

		switch r {
		case 119: // wを入力した場合、前の部屋に進む
			d.explore(front)
			return
		case 115: // sを入力した場合、後ろの部屋に戻る
			d.explore(back)
			return
		case 97: // aを入力した場合、左に方向転換
			d.player.Rotate(player.Left)
			return
		case 100: // dを入力した場合、右に方向転換
			d.player.Rotate(player.Right)
			return
		default: // 他の文字を入力した場合、再入力を促す
			fmt.Printf("\x1b[41m") // 背景色を赤に変更
			fmt.Printf("\n無効なキー操作です。行動を再入力してください。\n")
			fmt.Printf("\x1b[49m") // 背景色を戻す
			return
		}
	}
}

func (d *Dungeon) explore(dir room.Direction) {
	if d.currentRoom().HasDoor[dir] {
		// 進行方向にドアがあれば移動処理
		d.player.Move(dir)

		d.moveRoom = true // 移動検知

		system.System("clear")
		// ドアを開けるAA表示
		system.PrintFile("assets/rooms/openDoor")
		time.Sleep(400 * time.Millisecond)
	} else {
		// 進行方向にドアがない場合
		fmt.Printf("\x1b[41m") // 背景色を赤に変更
		fmt.Printf("\nその方向には進めません！\n")
		fmt.Printf("\x1b[49m") // 背景色を戻す
	}
}

func (d *Dungeon) currentRoom() *room.Room {
	return &d.rooms[d.player.RoomX][d.player.RoomY]
}
