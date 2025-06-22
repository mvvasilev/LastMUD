package game

import "github.com/google/uuid"

type Player struct {
	id uuid.UUID

	currentRoom *Room
}

func CreatePlayer(identity uuid.UUID, room *Room) *Player {
	return &Player{
		id:          identity,
		currentRoom: room,
	}
}

func (p *Player) Identity() string {
	return p.id.String()
}

func (p *Player) SetRoom(r *Room) {
	p.currentRoom = r
}
