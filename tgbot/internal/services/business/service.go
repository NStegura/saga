package business

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/tgbot/internal/clients/orders"
	"github.com/NStegura/saga/tgbot/internal/domain"
)

type Business struct {
	orderCli    IOrderCli
	paymentCli  IPaymentCli
	productsCli IProductCli

	logger *logrus.Logger
}

func New(orderCli IOrderCli, paymentCli IPaymentCli, productCli IProductCli, logger *logrus.Logger) *Business {
	return &Business{
		orderCli:    orderCli,
		paymentCli:  paymentCli,
		productsCli: productCli,
		logger:      logger,
	}
}

func (b *Business) GetProducts(ctx context.Context) ([]domain.Product, error) {
	products, err := b.productsCli.GetProducts(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	ps := make([]domain.Product, 0, len(products))
	for _, product := range products {
		ps = append(ps, domain.Product{
			ArticleID:   product.ArticleID,
			Category:    product.Category,
			Name:        product.Name,
			Description: product.Description,
			Count:       product.Count,
		})
	}
	return ps, nil
}

func (b *Business) GetProduct(ctx context.Context, articleID int64) (domain.Product, error) {
	product, err := b.productsCli.GetProduct(ctx, articleID)
	if err != nil {
		return domain.Product{}, fmt.Errorf("failed to get products: %w", err)
	}
	return domain.Product{
		ArticleID:   product.ArticleID,
		Category:    product.Category,
		Name:        product.Name,
		Description: product.Description,
		Count:       product.Count,
	}, nil
}

func (b *Business) CreateOrder(ctx context.Context, orderInfo domain.CreateOrderInfo) (int64, error) {
	oips := make([]orders.Product, 0, len(orderInfo.Products))
	for _, p := range orderInfo.Products {
		oips = append(oips, orders.Product{
			ArticleID: p.ArticleID,
			Count:     p.Count,
		})
	}
	orderID, err := b.orderCli.CreateOrder(ctx, orderInfo.UserID, orderInfo.Description, oips)
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	return orderID, nil
}

func (b *Business) GetOrders(ctx context.Context, userID int64) ([]domain.Order, error) {
	ords, err := b.orderCli.GetOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	ordersOut := make([]domain.Order, 0, len(ords))
	for _, order := range ords {
		o, err := b.GetOrder(ctx, order.OrderId)
		if err != nil {
			return nil, fmt.Errorf("failed to get order: %w", err)
		}
		ordersOut = append(ordersOut, o)
	}
	return ordersOut, nil
}

func (b *Business) GetOrder(ctx context.Context, orderID int64) (domain.Order, error) {
	order, err := b.orderCli.GetOrder(ctx, orderID)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to get orders: %w", err)
	}

	orderProducts := make([]domain.Product, 0, len(order.OrderProducts))
	for _, op := range order.OrderProducts {
		product, err := b.productsCli.GetProduct(ctx, op.ArticleID)
		if err != nil {
			return domain.Order{}, fmt.Errorf("failed to get product: %w", err)
		}
		orderProducts = append(orderProducts, domain.Product{
			ArticleID:   op.ArticleID,
			Category:    product.Category,
			Name:        product.Name,
			Description: product.Description,
			Count:       op.Count,
		})
	}
	statuses, err := b.orderCli.GetOrderStatuses(ctx, order.OrderId)
	if err != nil {
		return domain.Order{}, fmt.Errorf("failed to get order statuses: %w", err)
	}
	orderStatuses := make([]domain.StatusInfo, 0, len(statuses))
	for _, status := range statuses {
		orderStatuses = append(orderStatuses, domain.StatusInfo{
			Status: status.Status,
			Time:   status.Time,
		})
	}
	return domain.Order{
		OrderID:       order.OrderId,
		Products:      orderProducts,
		Description:   order.Description,
		CurrentStatus: order.State,
		StatusHistory: orderStatuses,
	}, nil
}

func (b *Business) PayOrder(ctx context.Context, orderID int64, status bool) error {
	err := b.paymentCli.PayOrder(ctx, orderID, status)
	if err != nil {
		return fmt.Errorf("failed to pay order: %w", err)
	}
	return nil
}
