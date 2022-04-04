package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"runtime"

	configPkg "github.com/iman_task/go-service/config"
	"github.com/iman_task/go-service/domain/service"
	"github.com/iman_task/go-service/events"
	"github.com/iman_task/go-service/events/handlers"
	pb "github.com/iman_task/go-service/genproto/collect"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
	broker "github.com/iman_task/go-service/pkg/messagebroker"
)

func main() {
	runtime.GOMAXPROCS(2)

	// =========================================================================
	// Configurations loading...
	config := configPkg.Load()

	// =========================================================================
	// Logger
	logger := loggerPkg.New("debug", "go-service")
	defer func() {
		err := loggerPkg.Cleanup(logger)
		if err != nil {
			logger.Fatal("failed cleaning up logs: %v", loggerPkg.Error(err))
		}
	}()

	// =========================================================================
	// Postgres
	logger.Info("Postgresql configs",
		loggerPkg.String("host", config.PostgresHost),
		loggerPkg.Int("port", config.PostgresPort),
		loggerPkg.String("database", config.PostgresDatabase),
	)

	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDatabase,
		config.PostgresSSL)

	// db initialization
	connDb, err := sqlx.Connect("postgres", psqlString)
	if err != nil {
		logger.Error("postgres connect error", loggerPkg.Error(err))
		return
	}
	connDb.SetMaxOpenConns(60)

	// =========================================================================
	// Kafka

	// Publishers
	publishersMap := make(map[string]broker.Producer)

	postAddTopicPublisher := events.NewKafkaProducer(&config, logger, handlers.PostAddTopic)
	defer func() {
		err := postAddTopicPublisher.Stop()
		if err != nil {
			logger.Fatal("Error while publishing: %v", loggerPkg.Error(err))
		}
	}()

	publishersMap[handlers.PostAddTopic] = postAddTopicPublisher

	// Listeners
	postChangeTopicListener := events.NewKafkaConsumer(connDb, &config, logger, handlers.PostChangeTopic)
	go postChangeTopicListener.Start()

	// =========================================================================
	// CronJobs
	lowStockCronJob := events.NewCronJob(connDb, &config, logger, publishersMap)
	go lowStockCronJob.Start()

	// =========================================================================
	// gRPC server
	goService := service.NewGoService(connDb, logger, config, publishersMap)

	listen, err := net.Listen("tcp", config.RPCPort)
	if err != nil {
		logger.Fatal("error while listening: %v", loggerPkg.Error(err))
	}
	s := grpc.NewServer()

	pb.RegisterGoServiceServer(s, goService)
	reflection.Register(s)

	logger.Info("main: server running", loggerPkg.String("port", config.RPCPort))

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Error while listening: %v", loggerPkg.Error(err))
	}
}
