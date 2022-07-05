package room

import (
	"fmt"
	"jailbreak/pkg/system"
)

// 東西南北の方向の定義
type Direction int

const (
	North Direction = iota
	East
	South
	West
)

type Room struct {
	HasDoor   [4]bool // 各方向のドアの有無
	IsVisited bool    // プレイヤーが部屋に来たことがあるかどうか
	HasHint   bool    // ヒントの表示有無
	IsWana    bool    // 罠の有無
	IsGoal    bool    // ゴールかどうか
}

func (r *Room) Display(dir Direction) {
	system.System("clear")
	if r.HasHint {
		fmt.Println("       ＼                                                        ／")
		fmt.Println("         ＼                          \x1b[33m★\x1b[39m                        ／")
	} else {
		fmt.Println("       ＼                                                        ／")
		fmt.Println("         ＼                                                    ／")
	}

	left := false  // プレーヤーから見て左のドア
	front := false // プレーヤーから見て前のドア
	right := false // プレーヤーから見て右のドア

	// プレーヤーの向きから左前右のドア情報取得
	if dir == North {
		left = r.HasDoor[West]
		front = r.HasDoor[North]
		right = r.HasDoor[East]
	} else if dir == East {
		left = r.HasDoor[North]
		front = r.HasDoor[East]
		right = r.HasDoor[South]
	} else if dir == South {
		left = r.HasDoor[East]
		front = r.HasDoor[South]
		right = r.HasDoor[West]
	} else if dir == West {
		left = r.HasDoor[South]
		front = r.HasDoor[West]
		right = r.HasDoor[North]
	} else {
		fmt.Println("向き情報に異常あり")
	}

	// 取得したドア情報から表示するAAを決定
	if !left {
		if !front {
			if !right {
				system.PrintFile("assets/rooms/door0")
			} else {
				system.PrintFile("assets/rooms/doorR")
			}
		} else {
			if !right {
				system.PrintFile("assets/rooms/doorF")
			} else {
				system.PrintFile("assets/rooms/doorFR")
			}
		}
	} else {
		if !front {
			if !right {
				system.PrintFile("assets/rooms/doorL")
			} else {
				system.PrintFile("assets/rooms/doorLR")
			}
		} else {
			if !right {
				system.PrintFile("assets/rooms/doorLF")
			} else {
				system.PrintFile("assets/rooms/doorLFR")
			}
		}
	}
}
