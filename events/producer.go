package events

import (
	"context"
	"fmt"
	congifPkg "github.com/iman_task/go-service/config"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
	brokerPkg "github.com/iman_task/go-service/pkg/messagebroker"
	"github.com/segmentio/kafka-go"
	"time"
)

type KafkaProducer struct {
	kafkaWriter *kafka.Writer
	logger      loggerPkg.Logger
}

func NewKafkaProducer(conf *congifPkg.Config, logger loggerPkg.Logger, topic string) brokerPkg.Producer {
	connString := fmt.Sprintf("%s:%d", conf.KafkaHost, conf.KafkaPort)

	return &KafkaProducer{
		kafkaWriter: kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{connString},
			Topic:        topic,
			BatchTimeout: 10 * time.Millisecond,
		}),
		logger: logger,
	}
}

// Start ...
func (p *KafkaProducer) Start() error {
	return nil
}

// Stop ...
func (p *KafkaProducer) Stop() error {
	err := p.kafkaWriter.Close()
	if err != nil {
		return err
	}

	return nil
}

// Publish ...
func (p *KafkaProducer) Publish(key, body []byte, logBody string) error {
	message := kafka.Message{
		Key:   key,
		Value: body,
	}

	if err := p.kafkaWriter.WriteMessages(context.Background(), message); err != nil {
		return err
	}

	p.logger.Info("Message published(key/body): " + string(key) + "/" + logBody)

	return nil
}
