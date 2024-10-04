package events

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/payments/internal/repo/models"
	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type PushJob struct {
	Frequency   time.Duration
	rateLimit   int
	eventsLimit int

	producer Producer
	repo     Repository
	logger   *logrus.Logger
}

func New(
	frequency time.Duration,
	rateLimit int,
	eventsLimit int,
	producer Producer,
	repo Repository,
	logger *logrus.Logger) *PushJob {
	return &PushJob{
		Frequency:   frequency,
		rateLimit:   rateLimit,
		eventsLimit: eventsLimit,
		producer:    producer,
		repo:        repo,
		logger:      logger,
	}
}

func (j *PushJob) Run(ctx context.Context) error {
	tx, err := j.repo.OpenTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to open transaction, %w", err)
	}

	eventsToPushCh, err := j.getEventsToPush(ctx, tx)
	if err != nil {
		_ = j.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to get waiting events: %w", err)
	}

	pushedEventsCh := make(chan int64, len(eventsToPushCh))
	var wg sync.WaitGroup
	for w := 1; w <= j.rateLimit; w++ {
		wg.Add(1)
		j.pushAndGetEvent(&wg, eventsToPushCh, pushedEventsCh)
	}
	go func() {
		wg.Wait()
		close(eventsToPushCh)
	}()

	for eventID := range pushedEventsCh {
		j.logger.Debugf("pushed event_id to update %v", eventID)
		err = j.repo.UpdateOutboxEvents(ctx, tx, []int64{eventID})
		if err != nil {
			j.logger.Error(err)
			continue
		}
	}
	_ = j.repo.Commit(ctx, tx)
	return nil
}

func (j *PushJob) getEventsToPush(ctx context.Context, tx pgx.Tx) (chan models.OutboxEntry, error) {
	eventsToPush, err := j.repo.GetNotProcessedEvents(ctx, tx, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get events from db: %w", err)
	}

	eventsToPushCh := make(chan models.OutboxEntry, len(eventsToPush))
	go func() {
		defer close(eventsToPushCh)
		for i := 0; i < len(eventsToPush); i++ {
			j.logger.Debug("get event to push", eventsToPush[i].ID)
			eventsToPushCh <- eventsToPush[i]
		}
	}()
	return eventsToPushCh, nil
}

func (j *PushJob) pushAndGetEvent(
	wg *sync.WaitGroup,
	eventsToPushCh chan models.OutboxEntry,
	pushedEventsCh chan int64) {
	go func() {
		defer wg.Done()

		for event := range eventsToPushCh {
			err := j.producer.PushMsg(event.Payload)
			if err != nil {
				j.logger.Error(err)
			} else {
				pushedEventsCh <- event.ID
			}
		}
	}()
}
