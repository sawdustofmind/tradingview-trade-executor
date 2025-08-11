package repository

import (
	"github.com/uptrace/bun"

	"github.com/frenswifbenefits/myfren/internal/crypt"
)

type Repository struct {
	crypter *crypt.Crypter

	db *bun.DB
}

func NewRepository(db *bun.DB, crypter *crypt.Crypter) *Repository {
	return &Repository{
		db:      db,
		crypter: crypter,
	}
}
