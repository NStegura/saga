package event

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/payments/internal/services/event/models"
	"github.com/sirupsen/logrus"
	"time"
)

type Event struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *Event {
	return &Event{repo: repo, logger: logger}
}

func (e *Event) ReserveEvents(ctx context.Context, reserveTo time.Time) (events []models.Event, err error) {
	tx, err := e.repo.OpenTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = e.repo.Commit(ctx, tx)
	}()

	eventsToPush, err := e.repo.GetNotProcessedEvents(ctx, tx, 5)
	if err != nil {
		_ = e.repo.Rollback(ctx, tx)
		return nil, fmt.Errorf("failed to get events from db: %w", err)
	}

	eventIDs := make([]int64, 0, len(eventsToPush))
	for _, eventToPush := range eventsToPush {
		eventIDs = append(eventIDs, eventToPush.ID)
	}
	err = e.repo.UpdateReservedTimeEvents(ctx, tx, eventIDs, reserveTo)
	if err != nil {
		_ = e.repo.Rollback(ctx, tx)
		return nil, fmt.Errorf("failed to get events from db: %w", err)
	}

	for _, ev := range eventsToPush {
		events = append(events, models.Event{
			ID:      ev.ID,
			Payload: ev.Payload,
			Topic:   ev.Topic,
		})
	}

	return events, nil
}

func (e *Event) UpdateEventStatusToDone(ctx context.Context, ID int64) error {
	e.logger.Debugf("pushed event_id to update %v", ID)
	tx, err := e.repo.OpenTransaction(ctx)
	if err != nil {
		e.logger.Error(err)
	}
	defer func() {
		_ = e.repo.Commit(ctx, tx)
	}()

	err = e.repo.UpdateEventStatusToDone(ctx, tx, ID)
	if err != nil {
		return fmt.Errorf("failed to update event status")
	}
	return nil
}
