package sender

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/golibs/event/service/models"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type PushEventsWorker struct {
	Frequency   time.Duration
	RateLimit   int
	EventsLimit int
	Reserve     time.Duration

	Producer Producer
	Event    Event
	Logger   *logrus.Logger
}

func (j *PushEventsWorker) GetFrequency() time.Duration {
	return j.Frequency
}

func (j *PushEventsWorker) Run(ctx context.Context) error {
	eventsToPushCh, err := j.reserveAndGetEventsToPush(ctx)
	if err != nil {
		return fmt.Errorf("failed to get waiting events: %w", err)
	}
	pushedEventsCh := make(chan int64, len(eventsToPushCh))
	var wg sync.WaitGroup
	for w := 1; w <= j.RateLimit; w++ {
		wg.Add(1)
		j.pushEvents(&wg, eventsToPushCh, pushedEventsCh)
	}
	go func() {
		wg.Wait()
		close(pushedEventsCh)
	}()

	for eventID := range pushedEventsCh {
		j.Logger.Debugf("pushed event_id to update %v", eventID)
		err := j.Event.UpdateEventStatusToDone(ctx, eventID)
		if err != nil {
			j.Logger.Error(err)
			continue
		}
	}
	return nil
}

func (j *PushEventsWorker) reserveAndGetEventsToPush(ctx context.Context) (chan models.Event, error) {
	eventsToPush, err := j.Event.ReserveEvents(ctx, int64(j.EventsLimit), time.Now().Add(j.Reserve))
	if err != nil {
		return nil, fmt.Errorf("failed to reserve events: %w", err)
	}

	eventsToPushCh := make(chan models.Event, len(eventsToPush))
	go func() {
		defer close(eventsToPushCh)
		for i := 0; i < len(eventsToPush); i++ {
			j.Logger.Debug("get event to push", eventsToPush[i].ID)
			eventsToPushCh <- eventsToPush[i]
		}
	}()
	return eventsToPushCh, nil
}

func (j *PushEventsWorker) pushEvents(
	wg *sync.WaitGroup,
	eventsToPushCh chan models.Event,
	pushedEventsCh chan int64) {
	go func() {
		defer wg.Done()

		for event := range eventsToPushCh {
			err := j.Producer.PushMsg(event.Payload, event.Topic)
			if err != nil {
				j.Logger.Error(err)
			} else {
				pushedEventsCh <- event.ID
			}
		}
	}()
}
