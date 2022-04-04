package handlers

import (
	config "github.com/iman_task/go-service/config"
	logger "github.com/iman_task/go-service/pkg/logger"
	broker "github.com/iman_task/go-service/pkg/messagebroker"
	"github.com/iman_task/go-service/storage"
)

const (
	PostAddTopic    = "post.add"
	PostChangeTopic = "post.change"
)

type EventHandler struct {
	conf      *config.Config
	storage   storage.Storage
	logger    logger.Logger
	publisher map[string]broker.Producer
}

func NewEventHandler(storage storage.Storage, logger logger.Logger, conf config.Config, publisher map[string]broker.Producer) *EventHandler {
	return &EventHandler{
		storage:   storage,
		conf:      &conf,
		logger:    logger,
		publisher: publisher,
	}
}

func (e *EventHandler) Handle(topic string, value []byte) (msg string, err error) {
	switch topic {
	case PostChangeTopic:
		msg, err = e.UpdatePost(value)
		if err != nil {
			return "", err
		}
	}

	return msg, nil
}
