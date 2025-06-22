package game

import "github.com/google/uuid"

type RoomPlayer interface {
	Identity() uuid.UUID
	SetRoom(room *Room)
}

type Room struct {
	world *World

	North *Room
	South *Room
	East  *Room
	West  *Room

	Name        string
	Description string

	players map[uuid.UUID]RoomPlayer
}

func CreateRoom(world *World, name, description string) *Room {
	return &Room{
		world:       world,
		Name:        name,
		Description: description,
		players:     map[uuid.UUID]RoomPlayer{},
	}
}

func (r *Room) PlayerJoinRoom(player RoomPlayer) (err error) {
	r.players[player.Identity()] = player

	return
}

func (r *Room) PlayerLeaveRoom(player RoomPlayer) (err error) {
	delete(r.players, player.Identity())

	return
}

func (r *Room) Players() map[uuid.UUID]RoomPlayer {
	return r.players
}

func (r *Room) World() *World {
	return r.world
}
