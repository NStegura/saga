package state

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/orders/internal/errs"
	"github.com/NStegura/saga/orders/internal/storage/repo/state/models"
)

type SRepo struct {
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(logger *logrus.Logger, pool *pgxpool.Pool) SRepo {
	return SRepo{logger: logger, pool: pool}
}

func (r *SRepo) CreateState(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	state models.OrderStateStatus,
) (stateID int64, err error) {
	const query = `
		INSERT INTO "state" (order_id, state) 
		VALUES ($1, $2) 
		RETURNING id
	`

	err = tx.QueryRow(ctx, query,
		orderID,
		state,
	).Scan(&orderID)
	if err != nil {
		return orderID, fmt.Errorf("failed to create order: %w", err)
	}
	return
}

func (r *SRepo) GetLastStateByOrderID(ctx context.Context, orderID int64) (state models.OrderState, err error) {
	const query = `
		SELECT id, order_id, state, created_at
		FROM "state"
		WHERE order_id = $1
		ORDER BY created_at DESC
		LIMIT 1;
	`

	err = r.pool.QueryRow(ctx, query, orderID).Scan(
		&state.ID,
		&state.OrderID,
		&state.State,
		&state.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errs.ErrNotFound
			return
		}
		return state, fmt.Errorf("failed to get order: %w", err)
	}
	return
}

func (r *SRepo) GetStatesByOrderID(ctx context.Context, orderID int64) (states []models.OrderState, err error) {
	const query = `
		SELECT id, order_id, state, created_at
		FROM "state"
		WHERE order_id = $1
		ORDER BY created_at;
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var state models.OrderState
		err = rows.Scan(
			&state.ID,
			&state.OrderID,
			&state.State,
			&state.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row products: %w", err)
		}
		states = append(states, state)
	}
	return states, nil
}
