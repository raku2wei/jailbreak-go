package dungeon

import (
	"fmt"
	"jailbreak/internal/room"
	"jailbreak/pkg/system"
	"time"
)

// 罠発動
func (d *Dungeon) wanaActivate() {
	fmt.Printf("侵入者発見！強制ワープします。\n")
	time.Sleep(2 * time.Second)

	// プレイヤーをスタート地点に強制移動
	d.player.SetPosition(1, 1, room.North)
	system.System("clear")
	d.Display()
	fmt.Printf("罠にかかったようだ\n")
}
