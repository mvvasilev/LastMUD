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

func (p *Player) Identity() uuid.UUID {
	return p.id
}

func (p *Player) SetRoom(r *Room) {
	p.currentRoom = r
}

func (p *Player) CurrentRoom() *Room {
	return p.currentRoom
}
