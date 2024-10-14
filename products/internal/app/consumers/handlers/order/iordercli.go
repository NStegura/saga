package order

import "github.com/NStegura/saga/products/internal/clients/orders"

type OrderCli interface {
	GetProductsToReserve(OrderID int64) (ps []orders.Product, err error)
}
