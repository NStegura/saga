package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/NStegura/saga/payments/internal/custom_errors"
	"github.com/NStegura/saga/payments/internal/repo/models"
	"github.com/jackc/pgx/v5"
	"time"
)

func (db *DB) CreatePayment(ctx context.Context, tx pgx.Tx, orderID int64) (id int64, err error) {
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

func (db *DB) UpdatePaymentStatusByOrderID(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	status models.PaymentStatus,
) (err error) {
	var id int64
	const query = `
		UPDATE "payment"
		SET	status = $1, updated_at = $2
		WHERE "payment".order_id = $3
		RETURNING "payment".id;
	`

	err = tx.QueryRow(ctx, query,
		status,
		time.Now(),
	).Scan(&orderID)

	if err != nil {
		return fmt.Errorf("UpdatePaymentStatus failed, %w", err)
	}
	db.logger.Debugf("Update payment status %v, id, %v", status, id)
	return
}

func (db *DB) GetPaymentByOrderID(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	forUpdate bool,
) (payment models.Payment, err error) {
	var query string
	if forUpdate {
		query = `
		SELECT p.id, p.order_id, p.status, p.created_at, p.updated_at
		FROM "payment" p
		WHERE id = $1
		FOR UPDATE;
	`
	} else {
		query = `
		SELECT p.id, p.order_id, p.status, p.created_at, p.updated_at
		FROM "payment" p
		WHERE id = $1;
	`
	}

	err = tx.QueryRow(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Status,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = custom_errors.ErrNotFound
			return
		}
		return payment, fmt.Errorf("get order failed, %w", err)
	}

	return
}