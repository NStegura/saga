package payment

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	dbModels "github.com/NStegura/saga/payments/internal/storage/repo/payment/models"
)

type EventRepository interface {
	CreateEvent(
		ctx context.Context,
		tx pgx.Tx,
		topic string,
		payload json.RawMessage,
	) (err error)
}

type Repository interface {
	EventRepository
	CreatePayment(ctx context.Context, tx pgx.Tx, orderID int64) (id int64, err error)
	UpdatePaymentStatusByOrderID(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		status dbModels.PaymentStatus,
	) error
	GetCreatedPaymentByOrderIDForUpdate(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
	) (payment dbModels.Payment, err error)

	OpenTransaction(ctx context.Context) (tx pgx.Tx, err error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
