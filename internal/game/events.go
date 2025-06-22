package game

import (
	"time"

	"github.com/google/uuid"
)

type PlayerJoinEvent struct {
	connectionId uuid.UUID
}

func CreatePlayerJoinEvent(connId uuid.UUID) *PlayerJoinEvent {
	return &PlayerJoinEvent{
		connectionId: connId,
	}
}

func (pje *PlayerJoinEvent) Type() EventType {
	return PlayerJoin
}

func (pje *PlayerJoinEvent) Handle(game *LastMUDGame, delta time.Duration) {
	game.world.AddPlayerToDefaultRoom(CreatePlayer(pje.connectionId, nil))
	game.enqeueOutput(CreateOutput(pje.connectionId, []byte("Welcome to LastMUD\n")))
}

type PlayerLeaveEvent struct {
	connectionId uuid.UUID
}

func CreatePlayerLeaveEvent(connId uuid.UUID) *PlayerLeaveEvent {
	return &PlayerLeaveEvent{
		connectionId: connId,
	}
}

func (ple *PlayerLeaveEvent) Type() EventType {
	return PlayerJoin
}

func (ple *PlayerLeaveEvent) Handle(game *LastMUDGame, delta time.Duration) {
	game.world.RemovePlayerById(ple.connectionId.String())
}
