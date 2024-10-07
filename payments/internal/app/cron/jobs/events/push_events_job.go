package events

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/payments/internal/services/event/models"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type PushJob struct {
	frequency   time.Duration
	rateLimit   int
	eventsLimit int
	reserve     time.Duration

	producer Producer
	event    Event
	logger   *logrus.Logger
}

func New(
	frequency time.Duration,
	rateLimit int,
	eventsLimit int,
	reserve time.Duration,
	producer Producer,
	event Event,
	logger *logrus.Logger) *PushJob {
	return &PushJob{
		frequency:   frequency,
		rateLimit:   rateLimit,
		eventsLimit: eventsLimit,
		reserve:     reserve,
		producer:    producer,
		event:       event,
		logger:      logger,
	}
}

func (j *PushJob) GetFrequency() time.Duration {
	return j.frequency
}

func (j *PushJob) Run(ctx context.Context) error {
	eventsToPushCh, err := j.reserveAndGetEventsToPush(ctx)
	if err != nil {
		return fmt.Errorf("failed to get waiting events: %w", err)
	}
	pushedEventsCh := make(chan int64, len(eventsToPushCh))
	var wg sync.WaitGroup
	for w := 1; w <= j.rateLimit; w++ {
		wg.Add(1)
		j.pushEvents(&wg, eventsToPushCh, pushedEventsCh)
	}
	go func() {
		wg.Wait()
		close(eventsToPushCh)
	}()

	for eventID := range pushedEventsCh {
		j.logger.Debugf("pushed event_id to update %v", eventID)
		err := j.event.UpdateEventStatusToDone(ctx, eventID)
		if err != nil {
			j.logger.Error(err)
			continue
		}
	}
	return nil
}

func (j *PushJob) reserveAndGetEventsToPush(ctx context.Context) (chan models.Event, error) {
	eventsToPush, err := j.event.ReserveEvents(ctx, time.Now().Add(j.reserve))
	if err != nil {
		return nil, fmt.Errorf("failed to reserve events: %w", err)
	}

	eventsToPushCh := make(chan models.Event, len(eventsToPush))
	go func() {
		defer close(eventsToPushCh)
		for i := 0; i < len(eventsToPush); i++ {
			j.logger.Debug("get event to push", eventsToPush[i].ID)
			eventsToPushCh <- eventsToPush[i]
		}
	}()
	return eventsToPushCh, nil
}

func (j *PushJob) pushEvents(
	wg *sync.WaitGroup,
	eventsToPushCh chan models.Event,
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
