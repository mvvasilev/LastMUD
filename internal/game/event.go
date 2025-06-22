package game

import "time"

type EventType int

const (
	PlayerJoin EventType = iota
	PlayerCommand
	PlayerLeave
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
		return event
	default:
		return nil
	}
}

func (eb *EventBus) Push(event GameEvent) {
	eb.events <- event
}

func (eb *EventBus) close() {
	close(eb.events)
}
