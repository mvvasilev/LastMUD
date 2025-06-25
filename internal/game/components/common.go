package components

import "code.haedhutner.dev/mvv/LastMUD/internal/game/ecs"

const (
	TypeName ecs.ComponentType = iota
	TypeDescription
	TypePlayerState
	TypeInRoom
	TypeNeighbors
	TypeIsRoom
	TypeIsPlayer
)

type NameComponent struct {
	Name string
}

func (c NameComponent) Type() ecs.ComponentType {
	return TypeName
}

type DescriptionComponent struct {
	Description string
}

func (c DescriptionComponent) Type() ecs.ComponentType {
	return TypeDescription
}
