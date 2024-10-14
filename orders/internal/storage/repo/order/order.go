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

type OrderRepo struct {
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(logger *logrus.Logger, pool *pgxpool.Pool) OrderRepo {
	return OrderRepo{logger: logger, pool: pool}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, tx pgx.Tx, userID int64, description string) (orderID int64, err error) {
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

func (r *OrderRepo) CreateProductOrder(ctx context.Context, tx pgx.Tx, orderId int64, productID int64, count int64) (ID int, err error) {
	const query = `
		INSERT INTO "product" (order_id, product_id, count) 
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err = tx.QueryRow(ctx, query,
		orderId, productID, count,
	).Scan(&ID)
	if err != nil {
		return ID, fmt.Errorf("failed to create product-order: %w", err)
	}
	return
}

func (r *OrderRepo) GetOrder(ctx context.Context, orderID int64) (order models.Order, err error) {
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

func (r *OrderRepo) GetOrders(ctx context.Context, userID int64) (orders []models.Order, err error) {
	const query = `
		SELECT id, description
		FROM "order"
		WHERE user_id = $1;
	`
	rows, err := r.pool.Query(ctx, query, userID)
	defer rows.Close()

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

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

func (r *OrderRepo) GetProductsByOrderId(ctx context.Context, orderID int64) (orderProduct []models.OrderProduct, err error) {
	const query = `
		SELECT id, order_id, product_id, count
		FROM "product"
		WHERE order_id = $1;
	`
	rows, err := r.pool.Query(ctx, query, orderID)
	defer rows.Close()

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get products: %w", err)
	}

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
