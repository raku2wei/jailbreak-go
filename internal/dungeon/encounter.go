package dungeon

import (
	"math/rand"
)

const (
	// エンカウント率のデフォルト値：10%
	DefaultEncounterRate float32 = 10.0
	// エンカウント率の増加率：移動するごとに1.4倍
	EncounterIncreaseRate float32 = 1.4
)

// 敵とのエンカウント処理
func (d *Dungeon) checkEncounter() bool {
	randomNumber := rand.Intn(100) + 1
	return float32(randomNumber) <= d.encounterRate
}
