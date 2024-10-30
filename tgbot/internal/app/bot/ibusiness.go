package bot

import (
	"context"

	"github.com/NStegura/saga/tgbot/internal/domain"
)

type Business interface {
	GetProducts(ctx context.Context) ([]domain.Product, error)
	GetProduct(ctx context.Context, articleID int64) (domain.Product, error)
	CreateOrder(ctx context.Context, orderInfo domain.CreateOrderInfo) (int64, error)
	GetOrders(ctx context.Context, userID int64) ([]domain.Order, error)
	GetOrder(ctx context.Context, orderID int64) (domain.Order, error)
	PayOrder(ctx context.Context, orderID int64, status bool) error
}
