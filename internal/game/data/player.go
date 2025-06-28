package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
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
