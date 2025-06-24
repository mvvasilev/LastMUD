package game

import (
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/game/command"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
	"github.com/google/uuid"
)

type PlayerJoinEvent struct {
	connectionId uuid.UUID
}

func (game *LastMUDGame) CreatePlayerJoinEvent(connId uuid.UUID) *PlayerJoinEvent {
	return &PlayerJoinEvent{
		connectionId: connId,
	}
}

func (pje *PlayerJoinEvent) Type() EventType {
	return PlayerJoin
}

func (pje *PlayerJoinEvent) Handle(game *LastMUDGame, delta time.Duration) {
	p := CreateJoiningPlayer(pje.connectionId)
	game.world.AddPlayerToDefaultRoom(p)
	game.enqeueOutput(game.CreateOutput(p.Identity(), []byte("Welcome to LastMUD!")))
	game.enqeueOutput(game.CreateOutput(p.Identity(), []byte("Please enter your name:")))
}

type PlayerLeaveEvent struct {
	connectionId uuid.UUID
}

func (game *LastMUDGame) CreatePlayerLeaveEvent(connId uuid.UUID) *PlayerLeaveEvent {
	return &PlayerLeaveEvent{
		connectionId: connId,
	}
}

func (ple *PlayerLeaveEvent) Type() EventType {
	return PlayerLeave
}

func (ple *PlayerLeaveEvent) Handle(game *LastMUDGame, delta time.Duration) {
	game.world.RemovePlayerById(ple.connectionId)
}

type PlayerCommandEvent struct {
	connectionId uuid.UUID
	command      *command.CommandContext
}

func (game *LastMUDGame) CreatePlayerCommandEvent(connId uuid.UUID, cmdString string) (event *PlayerCommandEvent, err error) {
	cmdCtx, err := command.CreateCommandContext(game.commandRegistry(), cmdString)

	if err != nil {
		return nil, err
	}

	event = &PlayerCommandEvent{
		connectionId: connId,
		command:      cmdCtx,
	}

	return
}

func (pce *PlayerCommandEvent) Type() EventType {
	return PlayerCommand
}

func (pce *PlayerCommandEvent) Handle(game *LastMUDGame, delta time.Duration) {
	player := game.world.FindPlayerById(pce.connectionId)

	if player == nil {
		logging.Error("Unable to handle player command from player with id", pce.connectionId, ": Player does not exist")
		return
	}

	event := pce.parseCommandIntoEvent(game, player)
}

func (pce *PlayerCommandEvent) parseCommandIntoEvent(game *LastMUDGame, player *Player) GameEvent {
	switch pce.command.Command().Definition().Name() {
	case SayCommand:
		speech, err := pce.command.Command().Parameters()[0].AsString()

		if err != nil {
			logging.Error("Unable to handle player speech from player with id", pce.connectionId, ": Speech could not be parsed: ", err.Error())
			return nil
		}

		return game.CreatePlayerSayEvent(player, speech)
	}

	return nil
}

type PlayerSayEvent struct {
	player *Player
	speech string
}

func (game *LastMUDGame) CreatePlayerSayEvent(player *Player, speech string) *PlayerSayEvent {
	return &PlayerSayEvent{
		player: player,
		speech: speech,
	}
}

func (pse *PlayerSayEvent) Type() EventType {
	return PlayerSpeak
}

func (pse *PlayerSayEvent) Handle(game *LastMUDGame, delta time.Duration) {
	for _, p := range pse.player.CurrentRoom().Players() {
		game.enqeueOutput(game.CreateOutput(p.Identity(), []byte(pse.player.id.String()+" in "+pse.player.CurrentRoom().Name+": "+pse.speech)))
	}
}
