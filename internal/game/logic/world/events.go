package world

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"github.com/google/uuid"
)

func CreatePlayerConnectEvent(world *ecs.World, connectionId uuid.UUID) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventPlayerConnect})
	ecs.SetComponent(world, event, data.ConnectionIdComponent{ConnectionId: connectionId})
}

func CreatePlayerDisconnectEvent(world *ecs.World, connectionId uuid.UUID) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventPlayerDisconnect})
	ecs.SetComponent(world, event, data.ConnectionIdComponent{ConnectionId: connectionId})
}

func CreatePlayerCommandEvent(world *ecs.World, connectionId uuid.UUID, command string) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventPlayerCommand})
	ecs.SetComponent(world, event, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, event, data.CommandStringComponent{Command: command})
}

func CreateParseCommandEvent(world *ecs.World, command ecs.Entity) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventParseCommand})
	ecs.SetComponent(world, event, data.EntityComponent{Entity: command})
}

func CreateCommandExecutedEvent(world *ecs.World, command ecs.Entity) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventCommandExecuted})
	ecs.SetComponent(world, event, data.EntityComponent{Entity: command})
}
