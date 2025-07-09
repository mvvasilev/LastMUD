package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
	"github.com/google/uuid"
)

func CreateGameOutput(w *ecs.World, connectionId uuid.UUID, contents string) ecs.Entity {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(w, gameOutput, data.IsOutputComponent{})
	ecs.SetComponent(w, gameOutput, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(w, gameOutput, data.ContentsComponent{Contents: []byte(contents)})

	return gameOutput
}

func CreateClosingGameOutput(w *ecs.World, connectionId uuid.UUID, contents []byte) ecs.Entity {
	gameOutput := ecs.NewEntity()

	ecs.SetComponent(w, gameOutput, data.IsOutputComponent{})
	ecs.SetComponent(w, gameOutput, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(w, gameOutput, data.ContentsComponent{Contents: contents})
	ecs.SetComponent(w, gameOutput, data.CloseConnectionComponent{})

	return gameOutput
}
