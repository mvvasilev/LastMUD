package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"github.com/google/uuid"
)

type ContentsComponent struct {
	Contents []byte
}

func (cc ContentsComponent) Type() ecs.ComponentType {
	return TypeContents
}

type CloseConnectionComponent struct{}

func (cc CloseConnectionComponent) Type() ecs.ComponentType {
	return TypeCloseConnection
}

func CreateGameOutput(world *ecs.World, connectionId uuid.UUID, contents []byte, shouldClose bool) {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(world, gameOutput, ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, gameOutput, ContentsComponent{Contents: contents})

	if shouldClose {
		ecs.SetComponent(world, gameOutput, CloseConnectionComponent{})
	}
}
