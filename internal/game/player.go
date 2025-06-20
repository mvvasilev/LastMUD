package game

type Player struct {
	GameObject
	Name
	Description
	Position
	Velocity
}

func CreatePlayer(name, description string, x, y int) *Player {
	return &Player{
		GameObject:  CreateGameObject(),
		Name:        WithName(name),
		Description: WithDescription(description),
		Position:    WithPosition(x, y),
		Velocity:    WithVelocity(0, 0),
	}
}
