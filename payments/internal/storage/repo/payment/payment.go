package payment

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/payments/internal/errs"
	"github.com/NStegura/saga/payments/internal/storage/repo/payment/models"
)

type PaymentRepo struct {
	logger *logrus.Logger
}

func New(logger *logrus.Logger) PaymentRepo {
	return PaymentRepo{logger: logger}
}

func (db *PaymentRepo) CreatePayment(ctx context.Context, tx pgx.Tx, orderID int64) (id int64, err error) {
	const query = `
		INSERT INTO "payment" (order_id, status) 
		VALUES ($1, $2) 
		RETURNING id;
	`

	err = tx.QueryRow(ctx, query,
		orderID,
		models.CREATED,
	).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("CreatePayment failed, %w", err)
	}
	db.logger.Debugf("Create payment, id, %v", id)
	return
}

func (db *PaymentRepo) UpdatePaymentStatusByOrderID(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	status models.PaymentStatus,
) (err error) {
	var ID int64
	const query = `
		UPDATE "payment"
		SET	status = $1, updated_at = $2
		WHERE "payment".order_id = $3
		RETURNING "payment".id;
	`

	err = tx.QueryRow(ctx, query,
		status,
		time.Now(),
		orderID,
	).Scan(&ID)

	if err != nil {
		return fmt.Errorf("UpdatePaymentStatus failed, %w", err)
	}
	db.logger.Debugf("Update payment status %v, id, %v", status, ID)
	return
}

func (db *PaymentRepo) GetCreatedPaymentByOrderIDForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
) (payment models.Payment, err error) {
	const query = `
		SELECT p.id, p.order_id, p.status, p.created_at, p.updated_at
		FROM "payment" p
		WHERE order_id = $1 and status = 'CREATED'
		FOR UPDATE;
	`
	err = tx.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errs.ErrNotFound
			return
		}
		return payment, fmt.Errorf("get order failed, %w", err)
	}

	return
}
