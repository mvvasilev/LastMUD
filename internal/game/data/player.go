package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
)

type PlayerState = byte

const (
	PlayerStateJoining PlayerState = iota
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

type InputComponent struct {
	Input rune
}

func (i InputComponent) Type() ecs.ComponentType {
	return TypeInput
}

type ExpectingDirectInput struct{}

func (e ExpectingDirectInput) Type() ecs.ComponentType {
	return TypeExpectingDirectInput
}

type ExpectingYNAnswer struct{}

func (e ExpectingYNAnswer) Type() ecs.ComponentType {
	return TypeExpectingYNAnswer
}
