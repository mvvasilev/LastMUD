package ecs

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

type NameComponent struct {
	Name string
}

type DescriptionComponent struct {
	Description string
}

type InRoomComponent struct {
	InRoom Entity
}

type NeighboringRoomsComponent struct {
	North, South, East, West Entity
}
