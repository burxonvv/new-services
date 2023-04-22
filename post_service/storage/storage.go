package storage

import (
	"github.com/new-york-services/post_service/storage/postgres"
	"github.com/new-york-services/post_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Post() repo.PostStorageI
}

type storagePg struct {
	db       *sqlx.DB
	postRepo repo.PostStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:       db,
		postRepo: postgres.NewPostRepo(db),
	}
}

func (s storagePg) Post() repo.PostStorageI {
	return s.postRepo
}
