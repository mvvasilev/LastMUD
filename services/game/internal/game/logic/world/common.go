package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
)

func QueryParent(e ecs.Entity) func(comp data.ParentComponent) bool {
	return func(comp data.ParentComponent) bool {
		return comp.Entity == e
	}
}

func QueryParentEvent(e ecs.Entity) func(comp data.ParentEventComponent) bool {
	return func(comp data.ParentEventComponent) bool {
		return comp.Event == e
	}
}
