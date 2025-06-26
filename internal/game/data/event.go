package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"github.com/google/uuid"
)

type EventType string

const (
	EventPlayerConnect    EventType = "PlayerConnect"
	EventPlayerDisconnect           = "PlayerDisconnect"
	EventPlayerCommand              = "PlayerCommand"
	EventPlayerSpeak                = "PlayerSpeak"
)

type EventComponent struct {
	EventType EventType
}

func (is EventComponent) Type() ecs.ComponentType {
	return TypeEvent
}

func CreatePlayerConnectEvent(world *ecs.World, connectionId uuid.UUID) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, EventComponent{EventType: EventPlayerConnect})
	ecs.SetComponent(world, event, ConnectionIdComponent{ConnectionId: connectionId})
}

func CreatePlayerDisconnectEvent(world *ecs.World, connectionId uuid.UUID) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, EventComponent{EventType: EventPlayerDisconnect})
	ecs.SetComponent(world, event, ConnectionIdComponent{ConnectionId: connectionId})
}

func CreatePlayerCommandEvent(world *ecs.World, connectionId uuid.UUID, command string) {
	event := ecs.NewEntity()

	ecs.SetComponent(world, event, EventComponent{EventType: EventPlayerCommand})
	ecs.SetComponent(world, event, ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, event, CommandStringComponent{Command: command})
}

// type PlayerJoinEvent struct {
// 	connectionId uuid.UUID
// }

// func (game *LastMUDGame) CreatePlayerJoinEvent(connId uuid.UUID) *PlayerJoinEvent {
// 	return &PlayerJoinEvent{
// 		connectionId: connId,
// 	}
// }

// func (pje *PlayerJoinEvent) Type() event.EventType {
// 	return PlayerJoin
// }

// func (pje *PlayerJoinEvent) Handle(game *LastMUDGame, delta time.Duration) {
// 	p, err := CreatePlayer(game.world.World, pje.connectionId, components.PlayerStateJoining)

// 	if err != nil {
// 		logging.Error("Unabled to create player: ", err)
// 	}

// 	game.enqeueOutput(game.CreateOutput(p.AsUUID(), []byte("Welcome to LastMUD!")))
// 	game.enqeueOutput(game.CreateOutput(p.AsUUID(), []byte("Please enter your name:")))
// }

// type PlayerLeaveEvent struct {
// 	connectionId uuid.UUID
// }

// func (game *LastMUDGame) CreatePlayerLeaveEvent(connId uuid.UUID) *PlayerLeaveEvent {
// 	return &PlayerLeaveEvent{
// 		connectionId: connId,
// 	}
// }

// func (ple *PlayerLeaveEvent) Type() event.EventType {
// 	return PlayerLeave
// }

// func (ple *PlayerLeaveEvent) Handle(game *LastMUDGame, delta time.Duration) {
// 	ecs.DeleteEntity(game.world.World, ecs.CreateEntity(ple.connectionId))
// }

// type PlayerCommandEvent struct {
// 	connectionId uuid.UUID
// 	command      *command.CommandContext
// }

// func (game *LastMUDGame) CreatePlayerCommandEvent(connId uuid.UUID, cmdString string) (event *PlayerCommandEvent, err error) {
// 	cmdCtx, err := command.CreateCommandContext(game.commandRegistry(), cmdString)

// 	if err != nil {
// 		return nil, err
// 	}

// 	event = &PlayerCommandEvent{
// 		connectionId: connId,
// 		command:      cmdCtx,
// 	}

// 	return
// }

// func (pce *PlayerCommandEvent) Type() event.EventType {
// 	return PlayerCommand
// }

// func (pce *PlayerCommandEvent) Handle(game *LastMUDGame, delta time.Duration) {
// 	if player == nil {
// 		logging.Error("Unable to handle player command from player with id", pce.connectionId, ": Player does not exist")
// 		return
// 	}

// 	event := pce.parseCommandIntoEvent(game, player)
// }

// func (pce *PlayerCommandEvent) parseCommandIntoEvent(game *LastMUDGame, player ecs.Entity) event.Event {
// 	switch pce.command.Command().Definition().Name() {
// 	case SayCommand:
// 		speech, err := pce.command.Command().Parameters()[0].AsString()

// 		if err != nil {
// 			logging.Error("Unable to handle player speech from player with id", pce.connectionId, ": Speech could not be parsed: ", err.Error())
// 			return nil
// 		}

// 		return game.CreatePlayerSayEvent(player, speech)
// 	}

// 	return nil
// }

// type PlayerSayEvent struct {
// 	player *Player
// 	speech string
// }

// func (game *LastMUDGame) CreatePlayerSayEvent(player *Player, speech string) *PlayerSayEvent {
// 	return &PlayerSayEvent{
// 		player: player,
// 		speech: speech,
// 	}
// }

// func (pse *PlayerSayEvent) Type() EventType {
// 	return PlayerSpeak
// }

// func (pse *PlayerSayEvent) Handle(game *LastMUDGame, delta time.Duration) {
// 	for _, p := range pse.player.CurrentRoom().Players() {
// 		game.enqeueOutput(game.CreateOutput(p.Identity(), []byte(pse.player.id.String()+" in "+pse.player.CurrentRoom().Name+": "+pse.speech)))
// 	}
// }
