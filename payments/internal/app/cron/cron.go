package cron

import (
	"context"
	"github.com/NStegura/saga/payments/internal/app/cron/jobs/events"
	"github.com/sirupsen/logrus"
	"time"
)

type Cron struct {
	pushEventJob events.PushJob
	logger       *logrus.Logger
}

func New(eventsJob events.PushJob, logger *logrus.Logger) *Cron {
	return &Cron{
		pushEventJob: eventsJob,
		logger:       logger,
	}
}

func (c *Cron) Start(ctx context.Context) error {
	pushEventsTimer := time.NewTicker(c.pushEventJob.Frequency)
	defer pushEventsTimer.Stop()
	i := 0
	for {
		select {
		case <-pushEventsTimer.C:
			i++
			c.logger.Infof("[PushEventsJob|%v] events push", i)
			if err := c.pushEventJob.Run(ctx); err != nil {
				c.logger.Error("PushEventsJob Run failed %v", err)
				continue
			}
		case <-ctx.Done():
			return nil
		}
	}
}
