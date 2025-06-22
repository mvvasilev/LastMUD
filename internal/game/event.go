package game

import (
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type EventType int

const (
	PlayerJoin EventType = iota
	PlayerCommand
	PlayerLeave

	PlayerSpeak
)

type GameEvent interface {
	Type() EventType
	Handle(game *LastMUDGame, delta time.Duration)
}

type EventBus struct {
	events chan GameEvent
}

func CreateEventBus() *EventBus {
	return &EventBus{
		events: make(chan GameEvent, 10),
	}
}

func (eb *EventBus) HasNext() bool {
	return len(eb.events) > 0
}

func (eb *EventBus) Pop() (event GameEvent) {
	select {
	case event := <-eb.events:
		logging.Info("Popped event of type ", event.Type(), ":", event)
		return event
	default:
		return nil
	}
}

func (eb *EventBus) Push(event GameEvent) {
	eb.events <- event
	logging.Info("Enqueued event of type ", event.Type(), ":", event)
}

func (eb *EventBus) close() {
	close(eb.events)
}
