package cron

import (
	"context"
	"time"

	"github.com/NStegura/saga/products/internal/app/cron/workers"
	"github.com/sirupsen/logrus"
)

type Cron struct {
	pushEventsWorker workers.Worker

	logger *logrus.Logger
}

func New(pushEventsWorker workers.Worker, logger *logrus.Logger) *Cron {
	return &Cron{
		pushEventsWorker: pushEventsWorker,
		logger:           logger,
	}
}

func (c *Cron) Start(ctx context.Context) error {
	pushEventsTimer := time.NewTicker(c.pushEventsWorker.GetFrequency())
	defer pushEventsTimer.Stop()
	i := 0
	for {
		select {
		case <-pushEventsTimer.C:
			i++
			c.logger.Infof("[ProductEventsWorker|%v] events push", i)
			if err := c.pushEventsWorker.Run(ctx); err != nil {
				c.logger.Error("ProductEventsWorker Run failed %v", err)
				continue
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Cron) Shutdown(_ context.Context) error {
	return nil
}

func (c *Cron) Name() string {
	return "cron tab"
}