package events

import (
	"context"
	dbModels "github.com/NStegura/saga/payments/internal/repo/models"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	GetNotProcessedEvents(
		ctx context.Context,
		tx pgx.Tx,
		limit int64,
	) (messages []dbModels.OutboxEntry, err error)
	UpdateOutboxEvents(ctx context.Context, tx pgx.Tx, messageIDs []int64) (err error)

	OpenTransaction(ctx context.Context) (tx pgx.Tx, err error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
