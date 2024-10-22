package product

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/products/internal/services/product/models"
	dbModels "github.com/NStegura/saga/products/internal/storage/repo/product/models"
)

const (
	topic = "inventory"
)

type Product struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *Product {
	return &Product{repo: repo, logger: logger}
}

func (p *Product) GetProducts(ctx context.Context) (ps []models.Product, err error) {
	products, err := p.repo.GetProducts(ctx)
	if err != nil {
		return ps, fmt.Errorf("failed to get products: %w", err)
	}
	outProducts := make([]models.Product, 0, len(products))
	for _, product := range products {
		outProducts = append(
			outProducts, models.Product{
				ProductID:   product.ID,
				Category:    product.Category,
				Name:        product.Name,
				Description: product.Description,
				Count:       product.Count,
			},
		)
	}
	return outProducts, nil
}

func (p *Product) GetProductInfo(ctx context.Context, productID int64) (pr models.Product, err error) {
	product, err := p.repo.GetProduct(ctx, productID)
	if err != nil {
		return pr, fmt.Errorf("failed to get product: %w", err)
	}
	return models.Product{
		ProductID:   product.ID,
		Category:    product.Category,
		Name:        product.Name,
		Description: product.Description,
		Count:       product.Count,
	}, nil
}

func (p *Product) ReserveProducts(
	ctx context.Context,
	orderID int64,
	reserves []models.Reserve,
) (err error) {
	var product dbModels.Product

	tx, err := p.repo.OpenTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = p.repo.Commit(ctx, tx)
	}()

	status := models.CREATED
	for _, r := range reserves {
		_, err = p.repo.CreateReserve(ctx, tx, orderID, r.ProductID, r.Count)
		if err != nil {
			_ = p.repo.Rollback(ctx, tx)
			return fmt.Errorf("failed to reserve product: %w", err)
		}

		product, err = p.repo.GetProductForUpdate(ctx, tx, r.ProductID)
		if err != nil {
			_ = p.repo.Rollback(ctx, tx)
			return fmt.Errorf("failed to get product for update: %w", err)
		}

		resCount := product.Count - r.Count

		if resCount < 0 {
			status = models.FAILED
			break
		}
		err = p.repo.UpdateProductCount(ctx, tx, product.ID, resCount)
		if err != nil {
			_ = p.repo.Rollback(ctx, tx)
			return fmt.Errorf("failed to update product count: %w", err)
		}
	}

	paymentMessage := models.ProductMessage{
		IKey:    uuid.New(),
		OrderID: orderID,
		Status:  status,
	}
	payload, err := json.Marshal(paymentMessage)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to marshal product message: %w", err)
	}
	err = p.repo.CreateEvent(ctx, tx, topic, payload)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to create event: %w", err)
	}

	return nil
}

func (p *Product) UpdateReserveStatus(ctx context.Context, orderID int64, status bool) (err error) {
	tx, err := p.repo.OpenTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to open transaction, %w", err)
	}
	defer func() {
		_ = p.repo.Commit(ctx, tx)
	}()

	if status == true {
		err = p.repo.UpdateReserveStatusByOrderID(ctx, tx, orderID, status)
		if err != nil {
			_ = p.repo.Rollback(ctx, tx)
			return fmt.Errorf("failed to reserve status by order id: %w", err)
		}
		return nil
	} else {
		reserves, err := p.repo.GetReservesByOrderIDForUpdate(ctx, tx, orderID)
		for _, reserve := range reserves {
			if err != nil {
				_ = p.repo.Rollback(ctx, tx)
				return fmt.Errorf("failed to get reserves: %w", err)
			}

			product, err := p.repo.GetProductForUpdate(ctx, tx, reserve.ProductID)
			if err != nil {
				_ = p.repo.Rollback(ctx, tx)
				return fmt.Errorf("failed to get product: %w", err)
			}

			resCount := product.Count + reserve.Count

			err = p.repo.UpdateProductCount(ctx, tx, product.ID, resCount)
			if err != nil {
				_ = p.repo.Rollback(ctx, tx)
				return fmt.Errorf("failed to update product count: %w", err)
			}

			err = p.repo.UpdateReserveStatusByID(ctx, tx, reserve.ID, false)
			if err != nil {
				_ = p.repo.Rollback(ctx, tx)
				return fmt.Errorf("failed to update reserve status: %w", err)
			}
		}
		return err
	}
}
