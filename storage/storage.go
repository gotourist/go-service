package storage

import (
	"database/sql"

	"github.com/iman_task/go-service/storage/postgres"
	"github.com/iman_task/go-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNoRows = sql.ErrNoRows
)

type Storage interface {
	Post() repo.PostStorage
	Collect() repo.CollectStorage
}
type storagePg struct {
	db          *sqlx.DB
	postRepo    repo.PostStorage
	collectRepo repo.CollectStorage
}

func NewStoragePg(db *sqlx.DB) Storage {
	return &storagePg{
		db:          db,
		postRepo:    postgres.NewPostRepo(db),
		collectRepo: postgres.NewCollectRepo(db),
	}
}

func (s storagePg) Post() repo.PostStorage {
	return s.postRepo
}

func (s storagePg) Collect() repo.CollectStorage {
	return s.collectRepo
}
