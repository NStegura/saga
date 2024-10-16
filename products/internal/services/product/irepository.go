package product

import (
	"context"
	"encoding/json"
	productRepo "github.com/NStegura/saga/products/internal/storage/repo/product/models"
	reserveRepo "github.com/NStegura/saga/products/internal/storage/repo/reserve/models"

	"github.com/jackc/pgx/v5"
)

type EventRepository interface {
	CreateEvent(
		ctx context.Context,
		tx pgx.Tx,
		topic string,
		payload json.RawMessage,
	) (err error)
}

type Repository interface {
	EventRepository
	GetProducts(context.Context) ([]productRepo.Product, error)
	GetProductForUpdate(ctx context.Context, tx pgx.Tx, productID int64) (product productRepo.Product, err error)
	GetProduct(ctx context.Context, productID int64) (product productRepo.Product, err error)
	UpdateProductCount(ctx context.Context, tx pgx.Tx, productID, count int64) (err error)

	CreateReserve(
		ctx context.Context,
		tx pgx.Tx,
		orderID,
		productID,
		count int64,
	) (id int64, err error)
	UpdateReserveStatusByOrderID(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
		status bool,
	) (err error)
	UpdateReserveStatusByID(
		ctx context.Context,
		tx pgx.Tx,
		ID int64,
		status bool,
	) (err error)
	GetReservesByOrderIDForUpdate(
		ctx context.Context,
		tx pgx.Tx,
		orderID int64,
	) (reserves []reserveRepo.Reserve, err error)

	OpenTransaction(ctx context.Context) (tx pgx.Tx, err error)
	Rollback(ctx context.Context, tx pgx.Tx) error
	Commit(ctx context.Context, tx pgx.Tx) error
}
