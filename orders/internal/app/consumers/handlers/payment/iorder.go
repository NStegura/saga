package payment

import (
	"context"

	"github.com/NStegura/saga/orders/internal/services/order/models"
)

// Order интерфейс для работы с бизнес слоем.
type Order interface {
	CreateOrderState(ctx context.Context, orderID int64, state models.OrderState) (err error)
}
