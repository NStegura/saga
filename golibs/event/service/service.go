package service

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/golibs/event/service/models"
	"github.com/sirupsen/logrus"
	"time"
)

type Event struct {
	Repo   Repository
	Logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *Event {
	return &Event{Repo: repo, Logger: logger}
}

func (e *Event) ReserveEvents(ctx context.Context, eventsLimit int64, reserveTo time.Time) (events []models.Event, err error) {
	tx, err := e.Repo.OpenTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = e.Repo.Commit(ctx, tx)
	}()

	eventsToPush, err := e.Repo.GetNotProcessedEvents(ctx, tx, eventsLimit)
	if err != nil {
		_ = e.Repo.Rollback(ctx, tx)
		return nil, fmt.Errorf("failed to get events from db: %w", err)
	}

	eventIDs := make([]int64, 0, len(eventsToPush))
	for _, eventToPush := range eventsToPush {
		eventIDs = append(eventIDs, eventToPush.ID)
	}
	err = e.Repo.UpdateReservedTimeEvents(ctx, tx, eventIDs, reserveTo)
	if err != nil {
		_ = e.Repo.Rollback(ctx, tx)
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
	e.Logger.Debugf("pushed event_id to update %v", ID)
	tx, err := e.Repo.OpenTransaction(ctx)
	if err != nil {
		e.Logger.Error(err)
	}
	defer func() {
		_ = e.Repo.Commit(ctx, tx)
	}()

	err = e.Repo.UpdateEventStatusToDone(ctx, tx, ID)
	if err != nil {
		return fmt.Errorf("failed to update event status")
	}
	return nil
}
