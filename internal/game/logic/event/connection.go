package event

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
)

func HandlePlayerConnect(w *ecs.World, event ecs.Entity) (err error) {
	connectionId, ok := ecs.GetComponent[data.ConnectionIdComponent](w, event)

	if !ok {
		return createEventHandlerError(data.EventPlayerConnect, "Event does not contain connectionId")
	}

	world.CreateJoiningPlayer(w, connectionId.ConnectionId)
	world.CreateGameOutput(w, connectionId.ConnectionId, "Welcome to LastMUD!")
	world.CreateGameOutput(w, connectionId.ConnectionId, "Before interacting with the game, you must either login or create a new account. Do so using the 'register' and 'login' command(s).")

	return
}

func HandlePlayerDisconnect(w *ecs.World, event ecs.Entity) (err error) {
	connectionId, ok := ecs.GetComponent[data.ConnectionIdComponent](w, event)

	if !ok {
		return createEventHandlerError(data.EventPlayerDisconnect, "Event does not contain connectionId")
	}

	playerEntity := ecs.QueryFirstEntityWithComponent(
		w,
		func(c data.ConnectionIdComponent) bool { return c.ConnectionId == connectionId.ConnectionId },
	)

	if playerEntity == ecs.NilEntity() {
		return createEventHandlerError(data.EventPlayerDisconnect, "Connection id cannot be associated with a player entity")
	}

	ecs.DeleteEntity(w, playerEntity)

	return
}
