package server

import (
	"context"

	"github.com/NStegura/saga/orders/internal/services/order/models"
)

// Order интерфейс для работы с бизнес слоем.
type Order interface {
	GetOrders(ctx context.Context, userID int64) ([]models.OrderInfo, error)
	GetOrder(ctx context.Context, orderID int64) (o models.Order, err error)
	GetOrderStates(ctx context.Context, orderID int64) (states []models.State, err error)
	CreateOrder(
		ctx context.Context,
		userID int64,
		description string,
		orderProducts []models.OrderProduct,
	) (orderID int64, err error)
}
