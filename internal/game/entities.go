package game

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/components"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/ecs"
	"github.com/google/uuid"
)

func CreatePlayer(world *ecs.World, id uuid.UUID, state components.PlayerState) (entity ecs.Entity, err error) {
	entity = ecs.CreateEntity(id)

	defaultRoom, err := ecs.GetResource[ecs.Entity](world, ResourceDefaultRoom)

	if err != nil {
		return
	}

	ecs.SetComponent(world, entity, components.PlayerStateComponent{State: state})
	ecs.SetComponent(world, entity, components.NameComponent{Name: id.String()})
	ecs.SetComponent(world, entity, components.InRoomComponent{Room: defaultRoom})
	ecs.SetComponent(world, entity, components.IsPlayerComponent{})

	return
}

func CreateRoom(
	world *ecs.World,
	name, description string,
	north, south, east, west ecs.Entity,
) ecs.Entity {
	entity := ecs.NewEntity()

	ecs.SetComponent(world, entity, components.IsRoomComponent{})
	ecs.SetComponent(world, entity, components.NameComponent{Name: name})
	ecs.SetComponent(world, entity, components.DescriptionComponent{Description: description})
	ecs.SetComponent(world, entity, components.NeighborsComponent{North: north, South: south, East: east, West: west})

	return entity
}
