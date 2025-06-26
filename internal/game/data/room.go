package data

import "code.haedhutner.dev/mvv/LastMUD/internal/ecs"

type IsRoomComponent struct {
}

func (c IsRoomComponent) Type() ecs.ComponentType {
	return TypeIsRoom
}

type NeighborsComponent struct {
	North, South, East, West ecs.Entity
}

func (c NeighborsComponent) Type() ecs.ComponentType {
	return TypeNeighbors
}

func CreateRoom(
	world *ecs.World,
	name, description string,
	north, south, east, west ecs.Entity,
) ecs.Entity {
	entity := ecs.NewEntity()

	ecs.SetComponent(world, entity, IsRoomComponent{})
	ecs.SetComponent(world, entity, NameComponent{Name: name})
	ecs.SetComponent(world, entity, DescriptionComponent{Description: description})
	ecs.SetComponent(world, entity, NeighborsComponent{North: north, South: south, East: east, West: west})

	return entity
}
