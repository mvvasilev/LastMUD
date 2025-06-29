package command

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
	"regexp"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
)

func HandleSay(w *ecs.World, _ time.Duration, player ecs.Entity, args data.ArgsMap) (err error) {
	playerRoom, ok := ecs.GetComponent[data.InRoomComponent](w, player)

	if !ok {
		return createCommandError("You aren't in a room!")
	}

	playerName, ok := ecs.GetComponent[data.NameComponent](w, player)

	if !ok {
		return createCommandError("You have no name!")
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
		world.SendMessageToPlayer(w, p, playerName.Name+": "+message)
	}

	return
}

func HandleQuit(w *ecs.World, _ time.Duration, player ecs.Entity, _ data.ArgsMap) (err error) {
	connId, _ := ecs.GetComponent[data.ConnectionIdComponent](w, player)

	world.CreateClosingGameOutput(w, connId.ConnectionId, []byte("Goodbye!"))

	return
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,24}$`)

func HandleRegister(w *ecs.World, delta time.Duration, player ecs.Entity, args data.ArgsMap) (err error) {
	accountName, err := arg[string](args, data.ArgAccountName)

	if err != nil {
		return err
	}

	if !usernameRegex.MatchString(accountName) {
		world.SendMessageToPlayer(w, player, "Registration: Username must only contain letters, numbers, dashes (-) and underscores (_), and be at most 24 characters in length.")
	}

	//accountPassword, err := arg[string](args, data.ArgAccountPassword)
	//
	//if err != nil {
	//	return err
	//}

	// TODO: validate username and password, encrypt password, etc.

	return
}
