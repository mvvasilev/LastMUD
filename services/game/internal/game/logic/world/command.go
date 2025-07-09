package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
)

func CreateTokenizedCommand(world *ecs.World, player ecs.Entity, commandString string, tokens []data.Token) ecs.Entity {
	command := ecs.NewEntity()

	ecs.SetComponent(world, command, data.ParentComponent{Entity: player})
	ecs.SetComponent(world, command, data.CommandStringComponent{Command: commandString})
	ecs.SetComponent(world, command, data.TokensComponent{Tokens: tokens})
	ecs.SetComponent(world, command, data.CommandStateComponent{State: data.CommandStateTokenized})

	return command
}
