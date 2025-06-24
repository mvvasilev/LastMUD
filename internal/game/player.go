package game

import "github.com/google/uuid"

type Player struct {
	id uuid.UUID

	state PlayerState

	currentRoom *Room
}

func CreateJoiningPlayer(identity uuid.UUID) *Player {
	return &Player{
		id:          identity,
		state:       PlayerStateJoining,
		currentRoom: nil,
	}
}

func CreatePlayer(identity uuid.UUID, state PlayerState, room *Room) *Player {
	return &Player{
		id:          identity,
		state:       state,
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
