package data

import (
	"code.haedhutner.dev/mvv/LastMUD/internal/ecs"
)

type IsOutputComponent struct{}

func (io IsOutputComponent) Type() ecs.ComponentType {
	return TypeIsOutput
}

type ContentsComponent struct {
	Contents []byte
}

func (cc ContentsComponent) Type() ecs.ComponentType {
	return TypeContents
}

type CloseConnectionComponent struct{}

func (cc CloseConnectionComponent) Type() ecs.ComponentType {
	return TypeCloseConnection
}
