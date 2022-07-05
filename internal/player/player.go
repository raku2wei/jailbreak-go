package player

import "jailbreak/internal/room"

type Rotation string

const (
	Left  Rotation = "left"
	Right Rotation = "right"
)

type Player struct {
	RoomX     int            // プレイヤーがいる部屋のX座標
	RoomY     int            // プレイヤーがいる部屋のY座標
	Direction room.Direction // 向いている方向
}

func NewPlayer() *Player {
	p := new(Player)
	return p
}

func (p *Player) SetPosition(x int, y int, dir room.Direction) {
	p.RoomX = x
	p.RoomY = y
	p.Direction = dir
}

func (p *Player) Move(dir room.Direction) {
	switch dir {
	case room.North:
		p.RoomY++
	case room.East:
		p.RoomX++
	case room.South:
		p.RoomY--
	case room.West:
		p.RoomX--
	}
}

func (p *Player) Rotate(r Rotation) {
	if r == Left {
		p.Direction--
		if p.Direction == -1 {
			p.Direction = 3
		}
	} else if r == Right {
		p.Direction++
		if p.Direction == 4 {
			p.Direction = 0
		}
	}
}
