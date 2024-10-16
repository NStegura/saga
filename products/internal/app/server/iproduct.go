package server

import (
	"context"

	"github.com/NStegura/saga/products/internal/services/product/models"
)

// Product интерфейс для работы с бизнес слоем.
type Product interface {
	GetProducts(context.Context) ([]models.Product, error)
	GetProductInfo(ctx context.Context, productID int64) (models.Product, error)
}
