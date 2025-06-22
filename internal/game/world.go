package game

type World struct {
	rooms       []*Room
	players     map[string]*Player
	defaultRoom *Room
}

func CreateWorld() *World {
	forest := CreateRoom("Forest", "A dense, misty forest stretches endlessly, its towering trees whispering secrets through rustling leaves. Sunbeams filter through the canopy, dappling the mossy ground with golden light.")
	cabin := CreateRoom("Wooden Cabin", "The cabin’s interior is cozy and rustic, with wooden beams overhead and a stone fireplace crackling warmly. A wool rug lies on creaky floorboards, and shelves brim with books, mugs, and old lanterns.")
	lake := CreateRoom("Ethermere Lake", "Ethermire Lake lies shrouded in mist, its dark, still waters reflecting a sky perpetually overcast. Whispers ride the wind, and strange lights flicker beneath the surface, never breaking it.")
	graveyard := CreateRoom("Graveyard", "An overgrown graveyard shrouded in fog, with cracked headstones and leaning statues. The wind sighs through dead trees, and unseen footsteps echo faintly among the mossy graves.")
	chapel := CreateRoom("Chapel of the Hollow Light", "This ruined chapel leans under ivy and age. Faint light filters through shattered stained glass, casting broken rainbows across dust-choked pews and a long-silent altar.")

	forest.North = cabin
	forest.South = graveyard
	forest.East = lake
	forest.West = chapel

	cabin.South = forest
	cabin.West = chapel
	cabin.East = lake

	chapel.North = cabin
	chapel.South = graveyard
	chapel.East = forest

	lake.West = forest
	lake.North = cabin
	lake.South = graveyard

	graveyard.North = forest
	graveyard.West = chapel
	graveyard.East = lake

	return &World{
		rooms: []*Room{
			forest,
			cabin,
			lake,
			graveyard,
			chapel,
		},
		defaultRoom: forest,
		players:     map[string]*Player{},
	}
}

func (w *World) AddPlayerToDefaultRoom(p *Player) {
	w.players[p.Identity()] = p
	w.defaultRoom.PlayerJoinRoom(p)
	p.SetRoom(w.defaultRoom)
}

func (w *World) RemovePlayerById(id string) {
	p, ok := w.players[id]

	if ok {
		p.currentRoom.PlayerLeaveRoom(p)
		delete(w.players, id)
		return
	}
}
