package game

import (
	"fmt"
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

func (et EventType) String() string {
	switch et {
	case PlayerCommand:
		return "PlayerCommand"
	case PlayerJoin:
		return "PlayerJoin"
	case PlayerLeave:
		return "PlayerLeave"
	case PlayerSpeak:
		return "PlayerSpeak"
	default:
		return "Unknown"
	}
}

type GameEvent interface {
	Type() EventType
	Handle(game *LastMUDGame, delta time.Duration)
}

func stringifyEvent(ev GameEvent) string {
	return ev.Type().String() + fmt.Sprintf(`%+v`, ev)
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
		logging.Debug("Popped event ", stringifyEvent(event))
		return event
	default:
		return nil
	}
}

func (eb *EventBus) Push(event GameEvent) {
	eb.events <- event
	logging.Debug("Enqueued event ", stringifyEvent(event))
}

func (eb *EventBus) close() {
	close(eb.events)
}
