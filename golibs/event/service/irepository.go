package service

import (
	"context"
	"github.com/NStegura/saga/golibs/event/repo/models"

	"github.com/jackc/pgx/v5"
	"time"
)

type Repository interface {
	GetNotProcessedEvents(
		ctx context.Context,
		tx pgx.Tx,
		limit int64,
	) (messages []models.EventEntry, err error)
	UpdateReservedTimeEvents(
		ctx context.Context,
		tx pgx.Tx,
		eventsIDs []int64,
		reservedTo time.Time,
	) (err error)
	UpdateEventStatusToDone(
		ctx context.Context,
		tx pgx.Tx,
		eventID int64,
	) (err error)

	OpenTransaction(ctx context.Context) (tx pgx.Tx, err error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
