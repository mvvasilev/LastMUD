package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"github.com/google/uuid"
)

type PlayerState = byte

const (
	PlayerStateJoining PlayerState = iota
	PlayerStateLoggingIn
	PlayerStateRegistering
	PlayerStatePlaying
	PlayerStateLeaving
)

type PlayerStateComponent struct {
	State PlayerState
}

func (c PlayerStateComponent) Type() ecs.ComponentType {
	return TypePlayerState
}

type InRoomComponent struct {
	Room ecs.Entity
}

func (c InRoomComponent) Type() ecs.ComponentType {
	return TypeInRoom
}

type IsPlayerComponent struct{}

func (c IsPlayerComponent) Type() ecs.ComponentType {
	return TypeIsPlayer
}

func CreatePlayer(world *ecs.World, id uuid.UUID, state PlayerState) (entity ecs.Entity, err error) {
	entity = ecs.NewEntity()

	defaultRoom, err := ecs.GetResource[ecs.Entity](world, ResourceDefaultRoom)

	if err != nil {
		return
	}

	ecs.SetComponent(world, entity, ConnectionIdComponent{ConnectionId: id})
	ecs.SetComponent(world, entity, PlayerStateComponent{State: state})
	ecs.SetComponent(world, entity, NameComponent{Name: id.String()})
	ecs.SetComponent(world, entity, InRoomComponent{Room: defaultRoom})
	ecs.SetComponent(world, entity, IsPlayerComponent{})

	return
}
