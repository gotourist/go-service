package repo

import (
	"github.com/iman_task/go-service/domain/entities"
)

type PostStorage interface {
	CreatePost(post *entities.Post) (*entities.Post, error)
	UpdatePost(post *entities.Post) error
}
