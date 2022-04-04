package events

import (
	"context"
	"fmt"
	"time"

	"github.com/iman_task/go-service/storage"
	"github.com/jmoiron/sqlx"

	"github.com/iman_task/go-service/config"
	"github.com/iman_task/go-service/events/handlers"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
	messageBrokerPkg "github.com/iman_task/go-service/pkg/messagebroker"
	"github.com/segmentio/kafka-go"
)

type KafkaConsumer struct {
	kafkaReader  *kafka.Reader
	eventHandler *handlers.EventHandler
	logger       loggerPkg.Logger
}

func NewKafkaConsumer(db *sqlx.DB, conf *config.Config, logger loggerPkg.Logger, topic string) messageBrokerPkg.Consumer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)

	return &KafkaConsumer{
		kafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{connString},
			GroupID:        "posts",
			Topic:          topic,
			MinBytes:       10e3, // 10KB
			MaxBytes:       10e6, // 10MB
			Partition:      0,
			CommitInterval: time.Second,
		}),
		eventHandler: handlers.NewEventHandler(storage.NewStoragePg(db), logger, *conf, nil),
		logger:       logger,
	}
}

func (k KafkaConsumer) Start() {
	fmt.Println(">>> Kafka consumer started.")
	for {
		m, err := k.kafkaReader.FetchMessage(context.Background())

		if err != nil {
			k.logger.Error("Error on consuming a message:", loggerPkg.Error(err))
			continue
		}

		msg, err := k.eventHandler.Handle(m.Topic, m.Value)

		if err != nil {
			k.logger.Error("failed to handle consumed topic:",
				loggerPkg.String("on topic", m.Topic), loggerPkg.Error(err))
		} else {
			k.logger.Info("Successfully consumed message",
				loggerPkg.String("on topic", m.Topic),
				loggerPkg.String("message", msg))
			k.kafkaReader.CommitMessages(context.Background(), m)
		}
	}
}
