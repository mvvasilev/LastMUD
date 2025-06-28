package world

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
)

func CreateRoom(
	world *ecs.World,
	name, description string,
	north, south, east, west ecs.Entity,
) ecs.Entity {
	entity := ecs.NewEntity()

	ecs.SetComponent(world, entity, data.IsRoomComponent{})
	ecs.SetComponent(world, entity, data.NameComponent{Name: name})
	ecs.SetComponent(world, entity, data.DescriptionComponent{Description: description})
	ecs.SetComponent(world, entity, data.NeighborsComponent{North: north, South: south, East: east, West: west})

	return entity
}
