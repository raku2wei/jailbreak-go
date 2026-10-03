package dungeon

import (
	"jailbreak/internal/console"
	"jailbreak/internal/room"
)

// ダンジョンマップ表示
func (d *Dungeon) PrintMap() {

	console.Printf("\n")

	for y := 4; y >= 1; y-- {
		// 1列目：北ドア
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].HasDoor[room.North] {
					console.Printf("____Ｄ____")
				} else {
					console.Printf("__________")
				}
			} else {
				console.Printf("          ")
			}
		}
		console.Printf("\n")

		// 2列目：見やすくするための空きスペース、ゴール
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].IsGoal {
					console.Printf("｜　G　｜")
				} else {
					console.Printf("｜      ｜")
				}
			} else {
				console.Printf("          ")
			}
		}
		console.Printf("\n")

		// 3列目：西ドア、プレイヤー、東ドア、罠
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].HasDoor[room.West] {
					console.Printf("Ｄ  ")
				} else {
					console.Printf("｜  ")
				}
				// 部屋にプレイヤーがいる場合
				if x == d.player.RoomX && y == d.player.RoomY {
					console.Printf("\x1b[36m")
					if d.player.Direction == room.North {
						console.Printf("↑ ")
					} else if d.player.Direction == room.East {
						console.Printf("→ ")
					} else if d.player.Direction == room.South {
						console.Printf("↓ ")
					} else if d.player.Direction == room.West {
						console.Printf("← ")
					} else {
						console.Printf("？")
					}
					console.Printf("\x1b[39m")
				} else {
					if d.rooms[x][y].IsWana {
						console.Printf("Ｗ")
					} else {
						console.Printf("  ")
					}
				}

				if d.rooms[x][y].HasDoor[room.East] {
					console.Printf("  Ｄ")
				} else {
					console.Printf("  ｜")
				}
			} else {
				console.Printf("          ")
			}
		}
		console.Printf("\n")

		// 4列目：見やすくするための空きスペース
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				console.Printf("｜      ｜")
			} else {
				console.Printf("          ")
			}
		}
		console.Printf("\n")

		// 5列目：南ドア
		for x := 1; x <= 6; x++ {
			if d.rooms[x][y].IsVisited {
				if d.rooms[x][y].HasDoor[room.South] {
					console.Printf("￣￣Ｄ￣￣")
				} else {
					console.Printf("￣￣￣￣￣")
				}
			} else {
				console.Printf("          ")
			}
		}
		console.Printf("\n")
	}
	// コンパス的な何か
	console.Printf("\x1b[36m")
	console.Printf("北\n")
	console.Printf("↑\n")
	console.Printf("　→東\n")
	console.Printf("\x1b[39m")
}
