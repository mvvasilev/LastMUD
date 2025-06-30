package command

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"regexp"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"

	"golang.org/x/crypto/bcrypt"
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

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,16}$`)
var passwordRegex = regexp.MustCompile(`^[a-zA-Z0-9!@#$%^&*()_+\-=\[\]{}|;:',.<>/?]{6,12}$`)

func HandleRegister(w *ecs.World, delta time.Duration, player ecs.Entity, args data.ArgsMap) (err error) {
	accountName, err := arg[string](args, data.ArgAccountName)

	if err != nil {
		return err
	}

	if !usernameRegex.MatchString(accountName) {
		world.SendMessageToPlayer(w, player, "Registration: Username must only contain letters, numbers, dashes (-) and underscores (_), and be at most 16 characters in length.")
		return
	}

	accountPassword, err := arg[string](args, data.ArgAccountPassword)

	if err != nil {
		return err
	}

	if !passwordRegex.MatchString(accountPassword) {
		world.SendMessageToPlayer(w, player, "Registration: Password must be between 6 and 12 characters in length")
		return
	}

	// TODO: Validate username doesn't exist already

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(accountPassword), bcrypt.DefaultCost)

	account := world.CreateAccount(w, accountName, encryptedPassword)

	ecs.SetComponent(w, player, data.AccountComponent{Account: account})

	world.SendMessageToPlayer(w, player, "Account created successfully! Welcome to LastMUD!")

	defaultRoom, err := ecs.GetResource[ecs.Entity](w, data.ResourceDefaultRoom)

	if err != nil {
		logging.Error("Unable to locate default room")
		world.SendMessageToPlayer(w, player, "Welcome to LastMUD! Your account was created, but you could not be joined to a room. Please try again later!")
		return
	}

	ecs.SetComponent(w, player, data.NameComponent{Name: accountName})
	ecs.SetComponent(w, player, data.InRoomComponent{Room: defaultRoom})
	ecs.SetComponent(w, player, data.PlayerStateComponent{State: data.PlayerStatePlaying})

	return
}
