package event

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/shared/log"
	"fmt"
	"time"
)

type Handler func(world *ecs.World, event ecs.Entity) (err error)

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

func CreateHandler(eventType data.EventType, handler Handler) ecs.SystemExecutor {
	return func(world *ecs.World, delta time.Duration) (err error) {
		events := ecs.QueryEntitiesWithComponent(world, eventTypeQuery(eventType))
		var processedEvents []ecs.Entity

		for event := range events {
			log.Debug("Handling event of type ", eventType)
			err = handler(world, event)

			if err != nil {
				log.Error(err)
			}

			processedEvents = append(processedEvents, event)
		}

		ecs.DeleteEntities(world, processedEvents...)

		return
	}
}
