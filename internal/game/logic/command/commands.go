package command

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
)

func HandleSay(w *ecs.World, _ time.Duration, player ecs.Entity, args data.ArgsMap) (err error) {
	playerRoom, ok := ecs.GetComponent[data.InRoomComponent](w, player)

	if !ok {
		return createCommandError("Player is not in any room!")
	}

	playerName, ok := ecs.GetComponent[data.NameComponent](w, player)

	if !ok {
		return createCommandError("Player has no name!")
	}

	allPlayersInRoom := ecs.QueryEntitiesWithComponent(w, func(comp data.InRoomComponent) bool {
		return comp.Room == playerRoom.Room
	})

	messageArg, ok := args[data.ArgMessageContent]

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
		connId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, p)

		world.CreateGameOutput(w, connId.ConnectionId, []byte(playerName.Name+": "+message))
	}

	return
}

func HandleQuit(w *ecs.World, _ time.Duration, player ecs.Entity, _ data.ArgsMap) (err error) {
	connId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, player)

	world.CreateClosingGameOutput(w, connId.ConnectionId, []byte("Goodbye!"))

	return
}

func HandleRegister(world *ecs.World, delta time.Duration, player ecs.Entity, args map[data.ArgName]data.Arg) (err error) {
	return
}
