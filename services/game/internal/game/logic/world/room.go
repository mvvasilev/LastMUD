package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
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

func MovePlayerToRoom(world *ecs.World, player, room ecs.Entity) {
	ecs.SetComponent(world, player, data.InRoomComponent{Room: room})
}
