package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
	"github.com/google/uuid"
)

func CreateJoiningPlayer(world *ecs.World, connectionId uuid.UUID) (entity ecs.Entity) {
	entity = ecs.NewEntity()

	ecs.SetComponent(world, entity, data.ConnectionIdComponent{ConnectionId: connectionId})
	ecs.SetComponent(world, entity, data.PlayerStateComponent{State: data.PlayerStateJoining})
	ecs.SetComponent(world, entity, data.NameComponent{Name: connectionId.String()})
	ecs.SetComponent(world, entity, data.IsPlayerComponent{})

	return
}

func SendMessageToPlayer(world *ecs.World, player ecs.Entity, message string) {
	connId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, player)

	if !ok {
		return
	}

	CreateGameOutput(world, connId.ConnectionId, message)
}

func SendDisconnectMessageToPlayer(world *ecs.World, player ecs.Entity, message string) {
	connId, ok := ecs.GetComponent[data.ConnectionIdComponent](world, player)

	if !ok {
		return
	}

	CreateClosingGameOutput(world, connId.ConnectionId, []byte(message))
}

func FindPlayerByConnectionId(w *ecs.World, connectionId uuid.UUID) (entity ecs.Entity) {
	player := ecs.NilEntity()

	for p := range ecs.IterateEntitiesWithComponent[data.IsPlayerComponent](w) {
		playerConnId, ok := ecs.GetComponent[data.ConnectionIdComponent](w, p)

		if ok && playerConnId.ConnectionId == connectionId {
			player = p
			break
		}
	}

	return player
}

func IsPlayerInDirectInputMode(w *ecs.World, player ecs.Entity) bool {
	_, ok := ecs.GetComponent[data.ExpectingDirectInput](w, player)
	return ok
}

func IsPlayerInYNAnswerMode(w *ecs.World, player ecs.Entity) bool {
	_, ok := ecs.GetComponent[data.ExpectingYNAnswer](w, player)
	return ok
}
