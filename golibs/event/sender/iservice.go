package sender

import (
	"context"
	"events/service/models"
	"time"
)

type Event interface {
	ReserveEvents(ctx context.Context, reserveTo time.Time) (events []models.Event, err error)
	UpdateEventStatusToDone(ctx context.Context, ID int64) error
}
