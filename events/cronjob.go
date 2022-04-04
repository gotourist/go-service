package events

import (
	"fmt"
	"github.com/iman_task/go-service/config"
	"github.com/iman_task/go-service/events/handlers"
	loggerPkg "github.com/iman_task/go-service/pkg/logger"
	messageBroker "github.com/iman_task/go-service/pkg/messagebroker"
	"github.com/iman_task/go-service/storage"
	"github.com/jmoiron/sqlx"
	"sync"
	"time"
)

type CronJob struct {
	eventHandler *handlers.EventHandler
	logger       loggerPkg.Logger
	storage      storage.Storage
}

func NewCronJob(db *sqlx.DB, conf *config.Config, logger loggerPkg.Logger, publisher map[string]messageBroker.Producer) *CronJob {
	return &CronJob{
		eventHandler: handlers.NewEventHandler(storage.NewStoragePg(db), logger, *conf, publisher),
		logger:       logger,
		storage:      storage.NewStoragePg(db),
	}
}

func (c *CronJob) Start() {
	fmt.Println(">>> Cron Job started.")
	ticker := time.NewTicker(time.Second)

	for range ticker.C {
		started, err := c.storage.Collect().CheckStarted()
		if err != nil {
			c.logger.Error("failed to check post collecting status from db", loggerPkg.Error(err))
			continue
		}

		if started {
			ticker.Stop()
			var wg sync.WaitGroup

			for i := 1; i <= 50; i++ {
				wg.Add(1)
				go c.eventHandler.Collect(&wg, i)
			}

			wg.Wait()

			err = c.storage.Collect().CollectPostsFinish()
			if err != nil {
				c.logger.Error("failed to finish collecting posts", loggerPkg.Error(err))
			}

			return
		}
	}
}
