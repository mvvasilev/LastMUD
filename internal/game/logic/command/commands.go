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

	message, err := arg[string](args, data.ArgMessageContent)

	if err != nil {
		return err
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

func HandleRegister(world *ecs.World, delta time.Duration, player ecs.Entity, args data.ArgsMap) (err error) {
	accountName, err := arg[string](args, data.ArgAccountName)

	if err != nil {
		return err
	}

	accountPassword, err := arg[string](args, data.ArgAccountPassword)

	if err != nil {
		return err
	}

	// TODO: validate username and password, encrypt password, etc.

	return
}
