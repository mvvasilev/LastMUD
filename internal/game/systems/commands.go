package systems

import (
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
)

func handleSayCommand(world *ecs.World, delta time.Duration, player ecs.Entity, args map[string]data.Arg) (err error) {
	playerRoom, ok := ecs.GetComponent[data.InRoomComponent](world, player)

	if !ok {
		return createCommandError("Player is not in any room!")
	}

	playerName, ok := ecs.GetComponent[data.NameComponent](world, player)

	if !ok {
		return createCommandError("Player has no name!")
	}

	allPlayersInRoom := ecs.QueryEntitiesWithComponent(world, func(comp data.InRoomComponent) bool {
		return comp.Room == playerRoom.Room
	})

	messageArg, ok := args["messageContent"]

	if !ok {
		return createCommandError("No message")
	}

	message, ok := messageArg.Value.(string)

	if !ok {
		return createCommandError("Can't interpret message as string")
	}

	if message == "" {
		return nil
	}

	for p := range allPlayersInRoom {
		connId, _ := ecs.GetComponent[data.ConnectionIdComponent](world, p)

		data.CreateGameOutput(world, connId.ConnectionId, []byte(playerName.Name+": "+message), false)
	}

	return
}

func handleQuitCommand(world *ecs.World, delta time.Duration, player ecs.Entity, _ map[string]data.Arg) (err error) {
	connId, _ := ecs.GetComponent[data.ConnectionIdComponent](world, player)

	data.CreateGameOutput(world, connId.ConnectionId, []byte("Goodbye!"), true)

	data.CreatePlayerDisconnectEvent(world, connId.ConnectionId)

	return
}
