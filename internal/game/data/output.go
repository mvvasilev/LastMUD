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

func CreateGameOutput(world *ecs.World, connectionId uuid.UUID, contents []byte) {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(world, gameOutput, ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, gameOutput, ContentsComponent{Contents: contents})
}
