package dungeon

import (
	"fmt"
	"jailbreak/internal/room"
)

// ダンジョンマップ表示
func (d *Dungeon) PrintMap() {

	fmt.Printf("\n")

	for y := 4; y >= 1; y-- {
		// 1列目：北ドア
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].HasDoor[room.North] {
					fmt.Printf("____Ｄ____")
				} else {
					fmt.Printf("__________")
				}
			} else {
				fmt.Printf("          ")
			}
		}
		fmt.Printf("\n")

		// 2列目：見やすくするための空きスペース、ゴール
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].IsGoal {
					fmt.Printf("｜　G　｜")
				} else {
					fmt.Printf("｜      ｜")
				}
			} else {
				fmt.Printf("          ")
			}
		}
		fmt.Printf("\n")

		// 3列目：西ドア、プレイヤー、東ドア、罠
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].HasDoor[room.West] {
					fmt.Printf("Ｄ  ")
				} else {
					fmt.Printf("｜  ")
				}
				// 部屋にプレイヤーがいる場合
				if x == d.player.RoomX && y == d.player.RoomY {
					fmt.Printf("\x1b[36m")
					if d.player.Direction == room.North {
						fmt.Printf("↑ ")
					} else if d.player.Direction == room.East {
						fmt.Printf("→ ")
					} else if d.player.Direction == room.South {
						fmt.Printf("↓ ")
					} else if d.player.Direction == room.West {
						fmt.Printf("← ")
					} else {
						fmt.Printf("？")
					}
					fmt.Printf("\x1b[39m")
				} else {
					if d.rooms[x][y].IsWana {
						fmt.Printf("Ｗ")
					} else {
						fmt.Printf("  ")
					}
				}

				if d.rooms[x][y].HasDoor[room.East] {
					fmt.Printf("  Ｄ")
				} else {
					fmt.Printf("  ｜")
				}
			} else {
				fmt.Printf("          ")
			}
		}
		fmt.Printf("\n")

		// 4列目：見やすくするための空きスペース
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				fmt.Printf("｜      ｜")
			} else {
				fmt.Printf("          ")
			}
		}
		fmt.Printf("\n")

		// 5列目：南ドア
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].HasDoor[room.South] {
					fmt.Printf("￣￣Ｄ￣￣")
				} else {
					fmt.Printf("￣￣￣￣￣")
				}
			} else {
				fmt.Printf("          ")
			}
		}
		fmt.Printf("\n")
	}
	// コンパス的な何か
	fmt.Printf("\x1b[36m")
	fmt.Printf("北\n")
	fmt.Printf("↑\n")
	fmt.Printf("　→東\n")
	fmt.Printf("\x1b[39m")
}
