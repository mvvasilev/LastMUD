package db

import "github.com/google/uuid"

type Identifier = uuid.UUID

type Entity interface {
	Id() Identifier
}

type Repository[T Entity] interface {
	Create(entity T) (rowsAffected int, err error)
	Delete(entity T) (rowsAffected int, err error)
	Update(entity T) (rowsAffected int, err error)
	FetchOne(id Identifier) (entity T, err error)
	FetchAll() (entities []T, err error)
}
