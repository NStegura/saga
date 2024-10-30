package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/orders/internal/errs"
	"github.com/NStegura/saga/orders/internal/storage/repo/order/models"
)

type ORepo struct {
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(logger *logrus.Logger, pool *pgxpool.Pool) ORepo {
	return ORepo{logger: logger, pool: pool}
}

func (r *ORepo) CreateOrder(
	ctx context.Context,
	tx pgx.Tx,
	userID int64,
	description string,
) (orderID int64, err error) {
	const query = `
		INSERT INTO "order" (user_id, description) 
		VALUES ($1, $2) 
		RETURNING id
	`

	err = tx.QueryRow(ctx, query,
		userID,
		description,
	).Scan(&orderID)
	if err != nil {
		return orderID, fmt.Errorf("failed to create order: %w", err)
	}
	return
}

func (r *ORepo) CreateProductOrder(
	ctx context.Context,
	tx pgx.Tx,
	orderID int64,
	productID int64,
	count int64,
) (id int, err error) {
	const query = `
		INSERT INTO "product" (order_id, product_id, count) 
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRow(ctx, query,
		orderID, productID, count,
	).Scan(&id)
	if err != nil {
		return id, fmt.Errorf("failed to create product-order: %w", err)
	}
	return
}

func (r *ORepo) GetOrder(ctx context.Context, orderID int64) (order models.Order, err error) {
	const query = `
		SELECT id, description
		FROM "order"
		WHERE id = $1;
	`

	err = r.pool.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.Description,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errs.ErrNotFound
			return
		}
		return order, fmt.Errorf("failed to get order: %w", err)
	}
	return
}

func (r *ORepo) GetOrders(ctx context.Context, userID int64) (orders []models.Order, err error) {
	const query = `
		SELECT id, description
		FROM "order"
		WHERE user_id = $1;
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		err = rows.Scan(
			&order.ID,
			&order.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row products: %w", err)
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *ORepo) GetProductsByOrderID(
	ctx context.Context,
	orderID int64,
) (orderProduct []models.OrderProduct, err error) {
	const query = `
		SELECT id, order_id, product_id, count
		FROM "product"
		WHERE order_id = $1;
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	defer rows.Close()

	var orderProducts []models.OrderProduct
	for rows.Next() {
		var orderProduct models.OrderProduct
		err = rows.Scan(
			&orderProduct.ID,
			&orderProduct.OrderID,
			&orderProduct.ProductID,
			&orderProduct.Count,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row products: %w", err)
		}
		orderProducts = append(orderProducts, orderProduct)
	}
	return orderProducts, nil
}
