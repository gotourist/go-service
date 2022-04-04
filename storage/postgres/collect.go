package postgres

import (
	"github.com/iman_task/go-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type collectRepo struct {
	db *sqlx.DB
}

func NewCollectRepo(db *sqlx.DB) repo.CollectStorage {
	return &collectRepo{
		db: db,
	}
}

func (c *collectRepo) CollectPostsStart() error {
	query := `
	UPDATE collect 
	SET "started"=true 
	WHERE "id" = 1`

	_, err := c.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (c *collectRepo) CollectPostsFinish() error {
	query := `
	UPDATE collect 
	SET "finished"=true 
	WHERE "id" = 1`

	_, err := c.db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (c *collectRepo) CheckFinished() (bool, error) {
	var finished bool

	query := `
		SELECT 
		"finished"
		FROM collect 
		WHERE "id" = 1`

	err := c.db.QueryRow(query).Scan(&finished)
	if err != nil {
		return false, err
	}

	return finished, nil
}

func (c *collectRepo) CheckStarted() (bool, error) {
	var started bool

	query := `
		SELECT 
		"started"
		FROM collect 
		WHERE "id" = 1`

	err := c.db.QueryRow(query).Scan(&started)
	if err != nil {
		return false, err
	}

	return started, nil
}
