package service

import (
	configPkg "github.com/iman_task/go-service/config"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
	broker "github.com/iman_task/go-service/pkg/messagebroker"
	"github.com/iman_task/go-service/storage"
	"github.com/jmoiron/sqlx"
)

type GoService struct {
	storage   storage.Storage
	logger    loggerPkg.Logger
	config    configPkg.Config
	publisher map[string]broker.Producer
}

func NewGoService(db *sqlx.DB, logger loggerPkg.Logger, config configPkg.Config, publisher map[string]broker.Producer) *GoService {
	return &GoService{
		storage:   storage.NewStoragePg(db),
		logger:    logger,
		config:    config,
		publisher: publisher,
	}
}
