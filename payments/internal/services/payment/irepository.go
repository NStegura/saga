package payment

import (
	"context"
	"encoding/json"
	dbModels "github.com/NStegura/saga/payments/internal/repo/models"
	"github.com/jackc/pgx/v5"
)

type Repository interface {
	CreatePayment(ctx context.Context, tx pgx.Tx, orderID int64) (id int64, err error)
	UpdatePaymentStatusByOrderID(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		status dbModels.PaymentStatus,
	) error
	GetPaymentByOrderID(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		forUpdate bool,
	) (payment dbModels.Payment, err error)
	CreateOutbox(
		ctx context.Context,
		tx pgx.Tx,
		payload json.RawMessage,
	) (err error)

	OpenTransaction(ctx context.Context) (tx pgx.Tx, err error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
