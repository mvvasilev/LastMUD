package ecs

type World struct {
	Players     []*Player
	Rooms       []*Room
	DefaultRoom *Room
}

func CreateWorld() *World {
	world := &World{
		Players: []*Player{},
		Rooms:   []*Room{},
	}

}
