package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/NStegura/saga/products/internal/errs"
	"github.com/NStegura/saga/products/internal/storage/repo/product/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type ProductRepo struct {
	pool *pgxpool.Pool

	logger *logrus.Logger
}

func New(logger *logrus.Logger, pool *pgxpool.Pool) ProductRepo {
	return ProductRepo{logger: logger, pool: pool}
}

func (r *ProductRepo) GetProducts(ctx context.Context) (ps []models.Product, err error) {
	const query = `
		SELECT id,
		       category,
		       name,
		       description,
		       count
		FROM product
		WHERE count > 0;
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		err = fmt.Errorf("failed to get products: %w", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var product models.Product
		err = rows.Scan(
			&product.ID,
			&product.Category,
			&product.Name,
			&product.Description,
			&product.Count,
		)
		ps = append(ps, product)
	}
	if errors.Is(err, pgx.ErrNoRows) {
		err = errs.ErrNotFound
		return
	}
	return
}

func (r *ProductRepo) GetProductForUpdate(
	ctx context.Context,
	tx pgx.Tx,
	productID int64,
) (product models.Product, err error) {
	const query = `
		SELECT id,
		       category,
		       name,
		       description,
		       count
		FROM product
		WHERE id = $1
		FOR UPDATE;
	`
	err = tx.QueryRow(ctx, query, productID).Scan(
		&product.ID,
		&product.Category,
		&product.Name,
		&product.Description,
		&product.Count,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err = errs.ErrNotFound
			return
		}
		return product, fmt.Errorf("get product for update failed, %w", err)
	}

	return
}

func (r *ProductRepo) UpdateProductCount(ctx context.Context, tx pgx.Tx, productID, count int64) (err error) {

	const query = `
		update product
		set	count = $2
		where id = $1;
	`

	cmd, err := tx.Exec(ctx, query,
		productID,
		count,
	)
	if err != nil {
		return fmt.Errorf("update product failed, %w", err)
	}

	if cmd.RowsAffected() == 0 {
		err = errs.ErrNotFound
		return
	}

	return
}
