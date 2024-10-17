package business

import (
	"context"

	"github.com/NStegura/saga/tgbot/internal/clients/products"
)

type IProductCli interface {
	GetProducts(ctx context.Context) ([]products.Product, error)
	GetProduct(ctx context.Context, articleID int64) (products.Product, error)
}
