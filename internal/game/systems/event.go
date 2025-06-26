package systems

import (
	"fmt"
	"time"

	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/logging"
)

type eventError struct {
	err string
}

func createEventError(v ...any) *eventError {
	return &eventError{
		err: fmt.Sprint(v...),
	}
}

func (e *eventError) Error() string {
	return e.err
}

func EventTypeQuery(eventType data.EventType) func(comp data.EventComponent) bool {
	return func(comp data.EventComponent) bool {
		return comp.EventType == eventType
	}
}

func CreateEventSystems() []*ecs.System {
	return []*ecs.System{
		ecs.CreateSystem("PlayerConnectEventHandlerSystem", 0, handlePlayerConnectEvents),
	}
}

func handlePlayerConnectEvents(world *ecs.World, delta time.Duration) (err error) {
	events := ecs.QueryEntitiesWithComponent(world, EventTypeQuery(data.EventPlayerConnect))
	processedEvents := []ecs.Entity{}

	for event := range events {
		err = handlePlayerConnectEvent(world, event)

		if err != nil {
			logging.Error("PlayerConnect Error: ", err)
		}

		processedEvents = append(processedEvents, event)
	}

	ecs.DeleteEntities(world, processedEvents...)

	return
}

func handlePlayerConnectEvent(world *ecs.World, entity ecs.Entity) (err error) {
	logging.Warn("Player connect")

	connectionId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, entity)

	if !ok {
		return createEventError("Event does not contain connectionId")
	}

	data.CreatePlayer(world, connectionId.ConnectionId, data.PlayerStateJoining)
	data.CreateGameOutput(world, connectionId.ConnectionId, []byte("Welcome to LastMUD!"))

	return
}
