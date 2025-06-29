package world

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"github.com/google/uuid"
)

func CreateJoiningPlayer(world *ecs.World, connectionId uuid.UUID) (entity ecs.Entity) {
	entity = ecs.NewEntity()

	ecs.SetComponent(world, entity, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, entity, data.PlayerStateComponent{State: data.PlayerStateJoining})
	ecs.SetComponent(world, entity, data.NameComponent{Name: connectionId.String()})
	ecs.SetComponent(world, entity, data.IsPlayerComponent{})
	ecs.SetComponent(world, entity, data.InputBufferComponent{InputBuffer: ""})

	return
}

func SendMessageToPlayer(world *ecs.World, player ecs.Entity, message string) {
	connId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, player)

	if !ok {
		return
	}

	CreateGameOutput(world, connId.ConnectionId, message)
}
