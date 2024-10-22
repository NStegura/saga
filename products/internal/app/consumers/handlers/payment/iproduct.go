package payment

import (
	"context"

	"github.com/NStegura/saga/products/internal/services/product/models"
)

// Product интерфейс для работы с бизнес слоем.
type Product interface {
	ReserveProducts(ctx context.Context, orderID int64, reserves []models.Reserve) (err error)
	UpdateReserveStatus(ctx context.Context, orderID int64, status bool) error
}
