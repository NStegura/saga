package orders

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/clients/orders/api"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/credentials/insecure"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn   *grpc.ClientConn
	client api.OrdersApiClient

	logger *logrus.Logger
}

func New(addr string, logger *logrus.Logger) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	return &Client{
		conn:   conn,
		client: api.NewOrdersApiClient(conn),
		logger: logger,
	}, nil
}

func (c *Client) GetOrder(ctx context.Context, orderID int64) (o Order, err error) {
	ctx = metadata.AppendToOutgoingContext(ctx,
		"sender", "tgbot",
		"when", time.Now().Format(time.RFC3339),
	)
	order, err := c.client.GetOrder(ctx, &api.OrderId{OrderId: orderID})
	if err != nil {
		return o, fmt.Errorf("failed to get order: %w", err)
	}
	orderProducts := make([]Product, 0, len(order.OrderProducts))
	for _, product := range order.OrderProducts {
		orderProducts = append(orderProducts, Product{
			ArticleID: product.ProductId,
			Count:     product.Count,
		})
	}

	return Order{
		OrderInfo: OrderInfo{
			OrderId:     order.OrderId,
			Description: order.Description,
			State:       order.State,
		},
		OrderProducts: orderProducts,
	}, nil
}

func (c *Client) GetOrders(ctx context.Context, userID int64) ([]OrderInfo, error) {
	ctx = metadata.AppendToOutgoingContext(ctx,
		"sender", "tgbot",
		"when", time.Now().Format(time.RFC3339),
	)
	orders, err := c.client.GetOrders(ctx, &api.UserId{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	ordersInfo := make([]OrderInfo, 0, len(orders.Orders))
	for _, orderInfo := range orders.Orders {
		ordersInfo = append(ordersInfo, OrderInfo{
			OrderId:     orderInfo.OrderId,
			Description: orderInfo.Description,
			State:       orderInfo.State,
		})
	}
	return ordersInfo, err
}

func (c *Client) GetOrderStatuses(ctx context.Context, orderID int64) ([]OrderStatus, error) {
	ctx = metadata.AppendToOutgoingContext(ctx,
		"sender", "tgbot",
		"when", time.Now().Format(time.RFC3339),
	)
	statuses, err := c.client.GetOrderStates(ctx, &api.OrderId{OrderId: orderID})
	if err != nil {
		return nil, fmt.Errorf("failed to get order statuses: %w", err)
	}
	orderStatuses := make([]OrderStatus, 0, len(statuses.OrderStates))
	for _, orderStatus := range statuses.OrderStates {
		orderStatuses = append(orderStatuses, OrderStatus{
			Status: orderStatus.State,
			Time:   orderStatus.Time.AsTime(),
		})
	}
	return orderStatuses, nil
}

func (c *Client) CreateOrder(ctx context.Context, userID int64, descr string, products []Product) (int64, error) {
	ctx = metadata.AppendToOutgoingContext(ctx,
		"sender", "tgbot",
		"when", time.Now().Format(time.RFC3339),
	)

	orderProducts := make([]*api.OrderProduct, 0, len(products))
	for _, p := range products {
		orderProducts = append(orderProducts, &api.OrderProduct{
			ProductId: p.ArticleID,
			Count:     p.Count,
		})
	}

	orderID, err := c.client.CreateOrder(ctx, &api.OrderIn{
		UserId:        userID,
		Description:   descr,
		OrderProducts: orderProducts,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to create order: %w", err)
	}
	return orderID.OrderId, nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
