package game

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/game/components"
	"code.haedhutner.dev/mvv/LastMUD/internal/game/ecs"
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

	forest := CreateRoom(
		gw.World,
		"Forest",
		"A dense, misty forest stretches endlessly, its towering trees whispering secrets through rustling leaves. Sunbeams filter through the canopy, dappling the mossy ground with golden light.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	ecs.SetResource(gw.World, ResourceDefaultRoom, forest)

	cabin := CreateRoom(
		gw.World,
		"Wooden Cabin",
		"The cabin’s interior is cozy and rustic, with wooden beams overhead and a stone fireplace crackling warmly. A wool rug lies on creaky floorboards, and shelves brim with books, mugs, and old lanterns.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	lake := CreateRoom(
		gw.World,
		"Ethermere Lake",
		"Ethermire Lake lies shrouded in mist, its dark, still waters reflecting a sky perpetually overcast. Whispers ride the wind, and strange lights flicker beneath the surface, never breaking it.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	graveyard := CreateRoom(
		gw.World,
		"Graveyard",
		"An overgrown graveyard shrouded in fog, with cracked headstones and leaning statues. The wind sighs through dead trees, and unseen footsteps echo faintly among the mossy graves.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	chapel := CreateRoom(
		gw.World,
		"Chapel of the Hollow Light",
		"This ruined chapel leans under ivy and age. Faint light filters through shattered stained glass, casting broken rainbows across dust-choked pews and a long-silent altar.",
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
		ecs.NilEntity(),
	)

	ecs.SetComponent(gw.World, forest, components.NeighborsComponent{
		North: cabin,
		South: graveyard,
		East:  lake,
		West:  chapel,
	})

	ecs.SetComponent(gw.World, cabin, components.NeighborsComponent{
		South: graveyard,
		West:  chapel,
		East:  lake,
	})

	ecs.SetComponent(gw.World, chapel, components.NeighborsComponent{
		North: cabin,
		South: graveyard,
		East:  forest,
	})

	ecs.SetComponent(gw.World, lake, components.NeighborsComponent{
		West:  forest,
		North: cabin,
		South: graveyard,
	})

	ecs.SetComponent(gw.World, graveyard, components.NeighborsComponent{
		North: forest,
		West:  chapel,
		East:  lake,
	})

	return
}

// type World struct {
// 	// rooms       []*Room
// 	// players     map[uuid.UUID]*Player
// 	// defaultRoom *Room

// 	entities   []Entity
// 	systems    []*System
// 	components map[ComponentType]any
// }

// func CreateWorld() (world *World) {
// 	world = &World{
// 		entities:   []Entity{},
// 		systems:    []*System{},
// 		components: map[ComponentType]any{},
// 	}
// 	// world = &World{
// 	// 	players: map[uuid.UUID]*Player{},
// 	// }

// 	// forest := CreateRoom(world, "Forest", "A dense, misty forest stretches endlessly, its towering trees whispering secrets through rustling leaves. Sunbeams filter through the canopy, dappling the mossy ground with golden light.")
// 	// cabin := CreateRoom(world, "Wooden Cabin", "The cabin’s interior is cozy and rustic, with wooden beams overhead and a stone fireplace crackling warmly. A wool rug lies on creaky floorboards, and shelves brim with books, mugs, and old lanterns.")
// 	// lake := CreateRoom(world, "Ethermere Lake", "Ethermire Lake lies shrouded in mist, its dark, still waters reflecting a sky perpetually overcast. Whispers ride the wind, and strange lights flicker beneath the surface, never breaking it.")
// 	// graveyard := CreateRoom(world, "Graveyard", "An overgrown graveyard shrouded in fog, with cracked headstones and leaning statues. The wind sighs through dead trees, and unseen footsteps echo faintly among the mossy graves.")
// 	// chapel := CreateRoom(world, "Chapel of the Hollow Light", "This ruined chapel leans under ivy and age. Faint light filters through shattered stained glass, casting broken rainbows across dust-choked pews and a long-silent altar.")

// 	// forest.North = cabin
// 	// forest.South = graveyard
// 	// forest.East = lake
// 	// forest.West = chapel

// 	// cabin.South = forest
// 	// cabin.West = chapel
// 	// cabin.East = lake

// 	// chapel.North = cabin
// 	// chapel.South = graveyard
// 	// chapel.East = forest

// 	// lake.West = forest
// 	// lake.North = cabin
// 	// lake.South = graveyard

// 	// graveyard.North = forest
// 	// graveyard.West = chapel
// 	// graveyard.East = lake

// 	// world.rooms = []*Room{
// 	// 	forest,
// 	// 	cabin,
// 	// 	lake,
// 	// 	graveyard,
// 	// 	chapel,
// 	// }

// 	// world.defaultRoom = forest

// 	return
// }

// func RegisterComponentType[T any](world *World, compType ComponentType) {
// 	world.components[compType] = CreateComponentStorage[any](compType)
// }

// func SetComponent[T](compType ComponentType, ent Entity, component any) {
// 	world.components[compType]
// }

// func (w *World) AddPlayerToDefaultRoom(p *Player) {
// 	w.players[p.Identity()] = p
// 	w.defaultRoom.PlayerJoinRoom(p)
// 	p.SetRoom(w.defaultRoom)
// }

// func (w *World) RemovePlayerById(id uuid.UUID) {
// 	p, ok := w.players[id]

// 	if ok {
// 		p.currentRoom.PlayerLeaveRoom(p)
// 		delete(w.players, id)
// 		return
// 	}
// }

// func (w *World) FindPlayerById(id uuid.UUID) *Player {
// 	p, ok := w.players[id]

// 	if ok {
// 		return p
// 	} else {
// 		return nil
// 	}
// }
