package ecs

import "github.com/google/uuid"

type Entity uuid.UUID

type Room struct {
	Entity
	NameComponent
	DescriptionComponent
	NeighboringRoomsComponent
}

type Player struct {
	Entity
	PlayerStateComponent
	NameComponent
	InRoomComponent
}
