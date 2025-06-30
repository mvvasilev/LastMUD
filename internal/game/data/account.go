package data

import "code.haedhutner.dev/mvv/LastMUD/internal/ecs"

type AccountComponent struct {
	Account ecs.Entity
}

func (ac AccountComponent) Type() ecs.ComponentType {
	return TypeAccount
}

type PasswordComponent struct {
	EncryptedPassword []byte
}

func (pc PasswordComponent) Type() ecs.ComponentType {
	return TypePassword
}
