package handlers

import (
	"github.com/iman_task/go-service/domain/entities"
	brokerpb "github.com/iman_task/go-service/genproto/broker/post"
	"github.com/iman_task/go-service/pkg/logger"
)

func (e *EventHandler) UpdatePost(data []byte) (msg string, err error) {
	var (
		postModel brokerpb.Post
		post      entities.Post
	)

	err = postModel.Unmarshal(data)
	if err != nil {
		return "", err
	}

	post = e.protoToPostModel(postModel)

	err = e.storage.Post().UpdatePost(&post)
	if err != nil {
		e.logger.Error("failed to update post in db", logger.Error(err))
		return "", err
	}

	return postModel.String(), nil
}

func (e *EventHandler) protoToPostModel(data brokerpb.Post) (post entities.Post) {
	post.Id = data.Id
	post.Title = data.Title
	post.Body = data.Body
	post.IsDeleted = data.IsDeleted

	return
}
