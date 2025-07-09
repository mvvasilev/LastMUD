package world

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/game/data"
)

func CreateAccount(world *ecs.World, username string, encryptedPassword []byte) ecs.Entity {
	account := ecs.NewEntity()

	ecs.SetComponent(world, account, data.NameComponent{Name: username})
	ecs.SetComponent(world, account, data.PasswordComponent{EncryptedPassword: encryptedPassword})

	return account
}
