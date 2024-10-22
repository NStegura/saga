package reserve

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/products/internal/errs"
	"github.com/NStegura/saga/products/internal/storage/repo/reserve/models"
)

type ReserveRepo struct {
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(logger *logrus.Logger, pool *pgxpool.Pool) ReserveRepo {
	return ReserveRepo{logger: logger, pool: pool}
}

func (r *ReserveRepo) CreateReserve(
	ctx context.Context,
	tx pgx.Tx,
	orderID,
	productID,
	count int64,
) (id int64, err error) {
	const query = `
		INSERT INTO "reserve" (product_id, order_id, count)
		VALUES ($1, $2, $3)
		RETURNING id;
	`

	err = tx.QueryRow(ctx, query,
		productID,
		orderID,
		count,
	).Scan(&id)

	if err != nil {
		return id, fmt.Errorf("CreateReserve failed, %w", err)
	}
	r.logger.Debugf("Create reserve, id: %v", id)
	return
}

func (r *ReserveRepo) UpdateReserveStatusByOrderID(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	status bool,
) (err error) {
	const query = `
		UPDATE "reserve"
		SET	pay_status = $2
		WHERE order_id = $1 and pay_status is NULL;
	`

	cmd, err := tx.Exec(ctx, query,
		orderID,
		status,
	)
	if err != nil {
		return fmt.Errorf("UpdateReverseStatusByOrderID failed, %w", err)
	}
	if cmd.RowsAffected() == 0 {
		err = errs.ErrNotFound
		return
	}

	return
}

func (r *ReserveRepo) UpdateReserveStatusByID(
	ctx context.Context,
	tx pgx.Tx,
	ID int64,
	status bool,
) (err error) {
	const query = `
		UPDATE "reserve"
		SET	pay_status = $2
		WHERE id = $1 and pay_status is NULL;
	`

	cmd, err := tx.Exec(ctx, query,
		ID,
		status,
	)
	if err != nil {
		return fmt.Errorf("UpdateReserveStatusByID failed, %w", err)
	}
	if cmd.RowsAffected() == 0 {
		err = errs.ErrNotFound
		return
	}

	return
}

func (r *ReserveRepo) GetReservesByOrderIDForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
) (reserves []models.Reserve, err error) {
	const query = `
		SELECT id, product_id, order_id, count, pay_status, saved_at
		FROM reserve
		WHERE order_id = $1 and pay_status is NULL
		FOR UPDATE;
	`
	rows, err := tx.Query(ctx, query, orderID)
	defer rows.Close()

	for rows.Next() {
		var reserve models.Reserve
		err = rows.Scan(
			&reserve.ID,
			&reserve.ProductID,
			&reserve.OrderID,
			&reserve.Count,
			&reserve.PayStatus,
			&reserve.SavedAt)
		reserves = append(reserves, reserve)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errs.ErrNotFound
			return nil, err
		}
		return reserves, fmt.Errorf("failed to get reserved ids: %w", err)
	}
	return reserves, nil
}
