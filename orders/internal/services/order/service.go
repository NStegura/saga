package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/NStegura/saga/orders/internal/services/order/models"
	dbStateModels "github.com/NStegura/saga/orders/internal/storage/repo/state/models"
)

const (
	topic = "order"
)

type Order struct {
	repo   Repository
	logger *logrus.Logger
}

func New(repo Repository, logger *logrus.Logger) *Order {
	return &Order{repo: repo, logger: logger}
}

func (p *Order) GetOrders(ctx context.Context, userID int64) ([]models.OrderInfo, error) {
	dbOrders, err := p.repo.GetOrders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	orders := make([]models.OrderInfo, 0, len(dbOrders))
	for _, dbOrder := range dbOrders {
		state, err := p.repo.GetLastStateByOrderID(ctx, dbOrder.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get state: %w", err)
		}

		orders = append(
			orders, models.OrderInfo{
				OrderID:     dbOrder.ID,
				Description: dbOrder.Description.String,
				State:       models.OrderState(state.State.String()),
			},
		)
	}
	return orders, nil
}

func (p *Order) GetOrder(ctx context.Context, orderID int64) (o models.Order, err error) {
	order, err := p.repo.GetOrder(ctx, orderID)
	if err != nil {
		return o, fmt.Errorf("failed to get order: %w", err)
	}

	state, err := p.repo.GetLastStateByOrderID(ctx, orderID)
	if err != nil {
		return o, fmt.Errorf("failed to get state: %w", err)
	}

	orderProducts, err := p.repo.GetProductsByOrderId(ctx, orderID)
	if err != nil {
		return o, fmt.Errorf("failed to get order products: %w", err)
	}

	ops := make([]models.OrderProduct, 0, len(orderProducts))
	for _, orderProduct := range orderProducts {
		ops = append(ops, models.OrderProduct{
			ProductID: orderProduct.ProductID,
			Count:     orderProduct.Count,
		})
	}
	return models.Order{
		OrderInfo: models.OrderInfo{
			OrderID:     order.ID,
			Description: order.Description.String,
			State:       models.OrderState(state.State.String()),
		},
		OrderProducts: ops,
	}, nil
}

func (p *Order) CreateOrder(
	ctx context.Context,
	userID int64,
	description string,
	orderProducts []models.OrderProduct,
) (orderID int64, err error) {
	tx, err := p.repo.OpenTransaction(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to open transaction: %w", err)
	}
	defer func() {
		_ = p.repo.Commit(ctx, tx)
	}()

	orderID, err = p.repo.CreateOrder(ctx, tx, userID, description)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return orderID, fmt.Errorf("failed to create order: %w", err)
	}

	for _, orderProduct := range orderProducts {
		_, err = p.repo.CreateProductOrder(ctx, tx, orderID, orderProduct.ProductID, orderProduct.Count)
		if err != nil {
			_ = p.repo.Rollback(ctx, tx)
			return orderID, fmt.Errorf("failed to create product-order: %w", err)
		}
	}

	orderMessage := models.OrderMessage{
		IKey:    uuid.New(),
		OrderID: orderID,
		Status:  models.CREATED,
	}
	payload, err := json.Marshal(orderMessage)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return orderID, fmt.Errorf("failed to marshal order message: %w", err)
	}
	err = p.repo.CreateEvent(ctx, tx, topic, payload)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return orderID, fmt.Errorf("failed to create event: %w", err)
	}

	_, err = p.repo.CreateState(ctx, tx, orderID, dbStateModels.ORDER_CREATED)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return orderID, fmt.Errorf("failed to update pay status: %w", err)
	}
	return orderID, nil
}

func (p *Order) CreateOrderState(ctx context.Context, orderID int64, state models.OrderState) (err error) {
	var dbState dbStateModels.OrderStateStatus

	tx, err := p.repo.OpenTransaction(ctx)
	if err != nil {
		return fmt.Errorf("failed to open transaction: %w", err)
	}
	defer func() {
		_ = p.repo.Commit(ctx, tx)
	}()

	err = dbState.Scan(string(state))
	if err != nil {
		return fmt.Errorf("failed to map order state status: %w", err)
	}
	_, err = p.repo.CreateState(ctx, tx, orderID, dbState)
	if err != nil {
		_ = p.repo.Rollback(ctx, tx)
		return fmt.Errorf("failed to update pay status: %w", err)
	}
	return nil
}
