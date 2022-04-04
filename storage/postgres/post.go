package postgres

import (
	"github.com/iman_task/go-service/domain/entities"
	"github.com/iman_task/go-service/storage/repo"
	"github.com/jmoiron/sqlx"
)

type postRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) repo.PostStorage {
	return &postRepo{
		db: db,
	}
}

func (p *postRepo) CreatePost(post *entities.Post) (*entities.Post, error) {
	var id int64

	tx, err := p.db.Beginx()
	if err != nil {
		return nil, err
	}

	stmt, err := tx.Preparex(`
		INSERT INTO 
		    post(
		         "title", 
		         "body"
		         ) 
		VALUES ($1, $2)
		RETURNING "id"`)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(
		post.Title,
		post.Body,
	).Scan(&id)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	query := `
		SELECT
		       "id", 
		       "title", 
		       "body", 
		       "is_deleted", 
		       "created_at" 
		FROM post 
		WHERE "id" = $1`

	err = p.db.Select(post, query, id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (p *postRepo) UpdatePost(post *entities.Post) error {
	tx, err := p.db.Beginx()
	if err != nil {
		return err
	}

	query := `
		UPDATE post 
		SET "title" = $1, 
		    "body" = $2,
		    "is_deleted" = $3
		WHERE "id" = $4`

	_, err = tx.Exec(
		query,
		post.Title,
		post.Body,
		post.IsDeleted,
		post.Id,
	)

	if err != nil {
		tx.Rollback()
		return err
	}

	return nil
}
