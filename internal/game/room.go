package game

type RoomPlayer interface {
	Identity() string
	SetRoom(room *Room)
}

type Room struct {
	North *Room
	South *Room
	East  *Room
	West  *Room

	Name        string
	Description string

	players map[string]RoomPlayer
}

func CreateRoom(name, description string) *Room {
	return &Room{
		Name:        name,
		Description: description,
		players:     map[string]RoomPlayer{},
	}
}

func (r *Room) PlayerJoinRoom(player RoomPlayer) (err error) {
	r.players[player.Identity()] = player

	return
}

func (r *Room) PlayerLeaveRoom(player RoomPlayer) (err error) {
	delete(r.players, player.Identity())

	return
}

func (r *Room) Players() map[string]RoomPlayer {
	return r.players
}
