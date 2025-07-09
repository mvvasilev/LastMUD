package data

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
)

type IsRoomComponent struct {
}

func (c IsRoomComponent) Type() ecs.ComponentType {
	return TypeIsRoom
}

type NeighborsComponent struct {
	North, South, East, West ecs.Entity
}

func (c NeighborsComponent) Type() ecs.ComponentType {
	return TypeNeighbors
}
