package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"github.com/google/uuid"
)

type EventType string

const (
	EventPlayerConnect    EventType = "PlayerConnect"
	EventPlayerDisconnect EventType = "PlayerDisconnect"
	EventPlayerCommand    EventType = "PlayerCommand"
	EventPlayerSpeak      EventType = "PlayerSpeak"
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
