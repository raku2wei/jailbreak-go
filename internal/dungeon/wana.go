package dungeon

import (
	"time"

	"jailbreak/internal/console"
	"jailbreak/internal/room"
)

// 罠発動
func (d *Dungeon) wanaActivate() {
	console.Printf("侵入者発見！強制ワープします。\n")
	time.Sleep(2 * time.Second)

	// プレイヤーをスタート地点に強制移動
	d.player.SetPosition(1, 1, room.North)
	console.Clear()
	d.Display()
	console.Printf("罠にかかったようだ\n")
}
