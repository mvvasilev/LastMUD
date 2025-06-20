package game

import "github.com/google/uuid"

type GameObject struct {
	uuid uuid.UUID
}

func CreateGameObject() GameObject {
	return GameObject{
		uuid: uuid.New(),
	}
}

type Position struct {
	x, y int
}

func WithPosition(x, y int) Position {
	return Position{x, y}
}

type Velocity struct {
	velX, velY int
}

func WithVelocity(velX, velY int) Velocity {
	return Velocity{velX, velY}
}

type Name struct {
	name string
}

func WithName(name string) Name {
	return Name{name}
}

type Description struct {
	description string
}

func WithDescription(description string) Description {
	return Description{description}
}
