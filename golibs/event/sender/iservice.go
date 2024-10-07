package sender

import (
	"context"
	"github.com/NStegura/saga/golibs/event/service/models"
	"time"
)

type Event interface {
	ReserveEvents(ctx context.Context, eventsLimit int64, reserveTo time.Time) (events []models.Event, err error)
	UpdateEventStatusToDone(ctx context.Context, ID int64) error
}
