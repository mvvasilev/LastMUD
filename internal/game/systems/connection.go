package systems

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

func handlePlayerConnectEvent(world *ecs.World, event ecs.Entity) (err error) {
	logging.Info("Player connect")

	connectionId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, event)

	if !ok {
		return createEventHandlerError(data.EventPlayerConnect, "Event does not contain connectionId")
	}

	data.CreatePlayer(world, connectionId.ConnectionId, data.PlayerStateJoining)
	data.CreateGameOutput(world, connectionId.ConnectionId, []byte("Welcome to LastMUD!"), false)

	return
}

func handlePlayerDisconnectEvent(world *ecs.World, event ecs.Entity) (err error) {
	logging.Info("Player disconnect")

	connectionId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, event)

	if !ok {
		return createEventHandlerError(data.EventPlayerDisconnect, "Event does not contain connectionId")
	}

	playerEntity := ecs.QueryFirstEntityWithComponent(
		world,
		func(c data.ConnectionIdComponent) bool { return c.ConnectionId == connectionId.ConnectionId },
	)

	if playerEntity == ecs.NilEntity() {
		return createEventHandlerError(data.EventPlayerDisconnect, "Connection id cannot be associated with a player entity")
	}

	ecs.DeleteEntity(world, playerEntity)

	return
}
