package data

import (
	"code.haedhutner.dev/mvv/LastMUD/services/game/internal/ecs"
)

type PasswordComponent struct {
	EncryptedPassword []byte
}

func (pc PasswordComponent) Type() ecs.ComponentType {
	return TypePassword
}
