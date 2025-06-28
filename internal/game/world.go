package game

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/data"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/logic/world"
)

type World struct {
	*ecs.World
}

func CreateGameWorld() (gw *World) {
	gw = &World{
		World: ecs.CreateWorld(),
	}

	defineRooms(gw.World)

	return
}

func defineRooms(w *ecs.World) {
	forest := world.CreateRoom(
		w,
		"Forest",
		"A dense, misty forest stretches endlessly, its towering trees whispering secrets through rustling leaves. Sunbeams filter through the canopy, dappling the mossy ground with golden light.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	ecs.SetResource(w, data.ResourceDefaultRoom, forest)

	cabin := world.CreateRoom(
		w,
		"Wooden Cabin",
		"The cabin’s interior is cozy and rustic, with wooden beams overhead and a stone fireplace crackling warmly. A wool rug lies on creaky floorboards, and shelves brim with books, mugs, and old lanterns.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	lake := world.CreateRoom(
		w,
		"Ethermere Lake",
		"Ethermire Lake lies shrouded in mist, its dark, still waters reflecting a sky perpetually overcast. Whispers ride the wind, and strange lights flicker beneath the surface, never breaking it.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	graveyard := world.CreateRoom(
		w,
		"Graveyard",
		"An overgrown graveyard shrouded in fog, with cracked headstones and leaning statues. The wind sighs through dead trees, and unseen footsteps echo faintly among the mossy graves.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	chapel := world.CreateRoom(
		w,
		"Chapel of the Hollow Light",
		"This ruined chapel leans under ivy and age. Faint light filters through shattered stained glass, casting broken rainbows across dust-choked pews and a long-silent altar.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	ecs.SetComponent(w, forest, data.NeighborsComponent{
		North: cabin,
		South: graveyard,
		East:  lake,
		West:  chapel,
	})

	ecs.SetComponent(w, cabin, data.NeighborsComponent{
		South: graveyard,
		West:  chapel,
		East:  lake,
	})

	ecs.SetComponent(w, chapel, data.NeighborsComponent{
		North: cabin,
		South: graveyard,
		East:  forest,
	})

	ecs.SetComponent(w, lake, data.NeighborsComponent{
		West:  forest,
		North: cabin,
		South: graveyard,
	})

	ecs.SetComponent(w, graveyard, data.NeighborsComponent{
		North: forest,
		West:  chapel,
		East:  lake,
	})
}
