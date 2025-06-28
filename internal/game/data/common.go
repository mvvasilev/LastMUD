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
	TypePlayer
	TypeCommandString
	TypeEntity
	TypeEvent
	TypeConnectionId
	TypeContents
	TypeCloseConnection

	TypeCommandTokens
	TypeCommandState
	TypeCommandArgs
	TypeCommand

	TypeAccount
	TypePassword
)

type Direction byte

const (
	DirectionNorth Direction = iota
	DirectionSouth
	DirectionEast
	DirectionWest
	DirectionUp
	DirectionDown
)

type PlayerComponent struct {
	Player ecs.Entity
}

func (p PlayerComponent) Type() ecs.ComponentType {
	return TypePlayer
}

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

type ConnectionIdComponent struct {
	ConnectionId uuid.UUID
}

func (cid ConnectionIdComponent) Type() ecs.ComponentType {
	return TypeConnectionId
}
