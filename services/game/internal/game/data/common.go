package data

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
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
	TypeInput
	TypeCommandString
	TypeParent
	TypeEvent
	TypeParentEvent
	TypeIsOutput
	TypeConnectionId
	TypeContents
	TypeCloseConnection

	TypeCommandTokens
	TypeCommandState
	TypeCommandArgs
	TypeCommand

	TypePassword
	TypeExpectingDirectInput
	TypeExpectingYNAnswer

	TypeFormFieldType
	TypeFormFieldValueString
	TypeFormFieldValueNumber
	TypeFormFieldValueDecimal
	TypeFormFieldValueYesNo
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

type ParentComponent struct {
	Entity ecs.Entity
}

func (e ParentComponent) Type() ecs.ComponentType {
	return TypeParent
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
