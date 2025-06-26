package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"github.com/google/uuid"
)

const (
	TypeName ecs.ComponentType = iota
	TypeDescription
	TypePlayerState
	TypeInRoom
	TypeNeighbors
	TypeIsRoom
	TypeIsPlayer
	TypeCommandString
	TypeEntity
	TypeEvent
	TypeConnectionId
	TypeContents
)

type EntityComponent struct {
	Entity ecs.Entity
}

func (e EntityComponent) Type() ecs.ComponentType {
	return TypeEntity
}

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

type CommandStringComponent struct {
	Command string
}

func (cs CommandStringComponent) Type() ecs.ComponentType {
	return TypeCommandString
}

type ConnectionIdComponent struct {
	ConnectionId uuid.UUID
}

func (cid ConnectionIdComponent) Type() ecs.ComponentType {
	return TypeConnectionId
}
