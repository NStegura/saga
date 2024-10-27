package order

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"

	"github.com/NStegura/saga/orders/internal/storage/repo/order/models"
	stateModels "github.com/NStegura/saga/orders/internal/storage/repo/state/models"
)

type EventRepository interface {
	CreateEvent(
		ctx context.Context,
		tx pgx.Tx,
		topic string,
		payload json.RawMessage,
	) (err error)
}

type orderRepository interface {
	GetOrder(ctx context.Context, orderID int64) (order models.Order, err error)
	GetOrders(ctx context.Context, userID int64) (orders []models.Order, err error)
	CreateOrder(ctx context.Context, tx pgx.Tx, userID int64, description string) (orderID int64, err error)
	CreateProductOrder(ctx context.Context, tx pgx.Tx, orderID int64, productID int64, count int64) (ID int, err error)
	GetProductsByOrderID(ctx context.Context, orderID int64) (orderProduct []models.OrderProduct, err error)
}

type orderStateRepository interface {
	CreateState(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		state stateModels.OrderStateStatus,
	) (stateID int64, err error)
	GetLastStateByOrderID(ctx context.Context, orderID int64) (state stateModels.OrderState, err error)
	GetStatesByOrderID(ctx context.Context, orderID int64) (states []stateModels.OrderState, err error)
}

type Repository interface {
	EventRepository
	orderRepository
	orderStateRepository
	OpenTransaction(ctx context.Context) (tx pgx.Tx, err error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
