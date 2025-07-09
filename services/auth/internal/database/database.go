package database

import (
	"gorm.io/gorm"
)

type AccountEntity struct {
	Username          string
	EncryptedPassword string

	gorm.Model
}
