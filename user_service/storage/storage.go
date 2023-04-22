package storage

import (
	"github.com/new-york-services/user_service/storage/postgres"
	"github.com/new-york-services/user_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	User() repo.UserStoreI
}

type storagePg struct {
	db       *sqlx.DB
	userRepo repo.UserStoreI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		userRepo: postgres.NewUserRepo(db),
	}
}

func (s storagePg) User() repo.UserStoreI {
	return s.userRepo
}
