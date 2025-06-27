package systems

import (
	"fmt"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type EventHandler func(world *ecs.World, event ecs.Entity) (err error)

type eventError struct {
	err string
}

func createEventHandlerError(eventType data.EventType, v ...any) *eventError {
	return &eventError{
		err: fmt.Sprint("Error handling ", string(eventType), ": ", v),
	}
}

func (e *eventError) Error() string {
	return e.err
}

func eventTypeQuery(eventType data.EventType) func(comp data.EventComponent) bool {
	return func(comp data.EventComponent) bool {
		return comp.EventType == eventType
	}
}

func CreateEventHandler(eventType data.EventType, handler EventHandler) ecs.SystemExecutor {
	return func(world *ecs.World, delta time.Duration) (err error) {
		events := ecs.QueryEntitiesWithComponent(world, eventTypeQuery(eventType))
		processedEvents := []ecs.Entity{}

		for event := range events {
			logging.Debug("Handling event of type ", eventType)
			err = handler(world, event)

			if err != nil {
				logging.Error(err)
			}

			processedEvents = append(processedEvents, event)
		}

		ecs.DeleteEntities(world, processedEvents...)

		return
	}
}

// func handlePlayerSayEvent(world *ecs.World, event ecs.Entity) (err error) {

// }
