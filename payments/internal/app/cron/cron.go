package cron

import (
	"context"
	"github.com/NStegura/saga/payments/internal/app/cron/jobs"
	"github.com/sirupsen/logrus"
	"time"
)

type Cron struct {
	eventJob jobs.Job

	logger *logrus.Logger
}

func New(eventsJob jobs.Job, logger *logrus.Logger) *Cron {
	return &Cron{
		eventJob: eventsJob,
		logger:   logger,
	}
}

func (c *Cron) Start(ctx context.Context) error {
	pushEventsTimer := time.NewTicker(c.eventJob.GetFrequency())
	defer pushEventsTimer.Stop()
	i := 0
	for {
		select {
		case <-pushEventsTimer.C:
			i++
			c.logger.Infof("[PaymentEventsJob|%v] events push", i)
			if err := c.eventJob.Run(ctx); err != nil {
				c.logger.Error("PaymentEventsJob Run failed %v", err)
				continue
			}
		case <-ctx.Done():
			return nil
		}
	}
}
