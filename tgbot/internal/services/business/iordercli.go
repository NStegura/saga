package business

import (
	"context"

	"github.com/NStegura/saga/tgbot/internal/clients/orders"
)

type IOrderCli interface {
	GetOrder(ctx context.Context, orderID int64) (o orders.Order, err error)
	GetOrders(ctx context.Context, userID int64) ([]orders.OrderInfo, error)
	GetOrderStatuses(ctx context.Context, orderID int64) ([]orders.OrderStatus, error)
	CreateOrder(ctx context.Context, userID int64, descr string, products []orders.Product) (int64, error)
}
