package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
)

type EventType string

const (
	EventPlayerConnect    EventType = "PlayerConnect"
	EventPlayerDisconnect EventType = "PlayerDisconnect"
	EventPlayerCommand    EventType = "PlayerCommand"
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
