package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
)

const (
	ResourceDefaultRoom ecs.Resource = "world:room:default"
)

type GameWorld struct {
	*ecs.World
}

func CreateGameWorld() (gw *GameWorld) {
	gw = &GameWorld{
		World: ecs.CreateWorld(),
	}

	defineRooms(gw.World)

	return
}

func defineRooms(world *ecs.World) {
	forest := CreateRoom(
		world,
		"Forest",
		"A dense, misty forest stretches endlessly, its towering trees whispering secrets through rustling leaves. Sunbeams filter through the canopy, dappling the mossy ground with golden light.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	ecs.SetResource(world, ResourceDefaultRoom, forest)

	cabin := CreateRoom(
		world,
		"Wooden Cabin",
		"The cabin’s interior is cozy and rustic, with wooden beams overhead and a stone fireplace crackling warmly. A wool rug lies on creaky floorboards, and shelves brim with books, mugs, and old lanterns.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	lake := CreateRoom(
		world,
		"Ethermere Lake",
		"Ethermire Lake lies shrouded in mist, its dark, still waters reflecting a sky perpetually overcast. Whispers ride the wind, and strange lights flicker beneath the surface, never breaking it.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	graveyard := CreateRoom(
		world,
		"Graveyard",
		"An overgrown graveyard shrouded in fog, with cracked headstones and leaning statues. The wind sighs through dead trees, and unseen footsteps echo faintly among the mossy graves.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	chapel := CreateRoom(
		world,
		"Chapel of the Hollow Light",
		"This ruined chapel leans under ivy and age. Faint light filters through shattered stained glass, casting broken rainbows across dust-choked pews and a long-silent altar.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	ecs.SetComponent(world, forest, NeighborsComponent{
		North: cabin,
		South: graveyard,
		East:  lake,
		West:  chapel,
	})

	ecs.SetComponent(world, cabin, NeighborsComponent{
		South: graveyard,
		West:  chapel,
		East:  lake,
	})

	ecs.SetComponent(world, chapel, NeighborsComponent{
		North: cabin,
		South: graveyard,
		East:  forest,
	})

	ecs.SetComponent(world, lake, NeighborsComponent{
		West:  forest,
		North: cabin,
		South: graveyard,
	})

	ecs.SetComponent(world, graveyard, NeighborsComponent{
		North: forest,
		West:  chapel,
		East:  lake,
	})
}
