package data

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
)

type EventType string

const (
	EventPlayerConnect    EventType = "PlayerConnect"
	EventPlayerDisconnect EventType = "PlayerDisconnect"
	EventPlayerInput      EventType = "PlayerInput"
	EventSubmitInput      EventType = "PlayerCommand"
	EventParseCommand     EventType = "ParseCommand"
	EventCommandExecuted  EventType = "CommandExecuted"
	EventPlayerSpeak      EventType = "PlayerSpeak"
)

type EventComponent struct {
	EventType EventType
}

func (is EventComponent) Type() ecs.ComponentType {
	return TypeEvent
}

type ParentEventComponent struct {
	Event ecs.Entity
}

func (c ParentEventComponent) Type() ecs.ComponentType {
	return TypeParentEvent
}
