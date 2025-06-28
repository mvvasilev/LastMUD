package world

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"github.com/google/uuid"
)

func CreateGameOutput(w *ecs.World, connectionId uuid.UUID, contents []byte) ecs.Entity {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(w, gameOutput, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(w, gameOutput, data.ContentsComponent{Contents: contents})

	return gameOutput
}

func CreateClosingGameOutput(w *ecs.World, connectionId uuid.UUID, contents []byte) ecs.Entity {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(w, gameOutput, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(w, gameOutput, data.ContentsComponent{Contents: contents})
	ecs.SetComponent(w, gameOutput, data.CloseConnectionComponent{})

	return gameOutput
}
