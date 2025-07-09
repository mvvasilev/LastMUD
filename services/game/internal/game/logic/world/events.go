package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
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

func CreatePlayerInputEvent(world *ecs.World, connectionId uuid.UUID, input rune) {
	player := ecs.QueryFirstEntityWithComponent[data.ConnectionIdComponent](world, func(comp data.ConnectionIdComponent) bool {
		return comp.ConnectionId == connectionId
	})

	if player == ecs.NilEntity() {
		log.Error("Trying to process input event for connection '", connectionId.String(), "' which does not have a corresponding player")
		return
	}

	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventPlayerInput})
	ecs.SetComponent(world, event, data.ParentComponent{Entity: player})
	ecs.SetComponent(world, event, data.InputComponent{Input: input})
}

func CreateSubmitInputEvent(world *ecs.World, connectionId uuid.UUID, command string) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventSubmitInput})
	ecs.SetComponent(world, event, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, event, data.CommandStringComponent{Command: command})
}

func CreateParseCommandEvent(world *ecs.World, command ecs.Entity) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventParseCommand})
	ecs.SetComponent(world, command, data.ParentEventComponent{Event: event})
}

func CreateCommandExecutedEvent(world *ecs.World, command ecs.Entity) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, data.EventComponent{EventType: data.EventCommandExecuted})
	ecs.SetComponent(world, command, data.ParentEventComponent{Event: event})
}
