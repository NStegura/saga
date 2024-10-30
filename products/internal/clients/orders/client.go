package orders

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/NStegura/saga/products/internal/clients/orders/api"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn   *grpc.ClientConn
	client api.OrdersApiClient
}

func New(addr string) (*Client, error) {
	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}
	return &Client{
		conn:   conn,
		client: api.NewOrdersApiClient(conn),
	}, nil
}

func (c *Client) GetProductsToReserve(orderID int64) (ps []Product, err error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"sender", "productService",
		"when", time.Now().Format(time.RFC3339),
	)
	order, err := c.client.GetOrder(ctx, &api.OrderId{OrderId: orderID})
	if err != nil {
		return ps, fmt.Errorf("failed to get order: %w", err)
	}
	for _, product := range order.OrderProducts {
		ps = append(ps, Product{
			ProductID: product.ProductId,
			Count:     product.Count,
		})
	}

	return ps, nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
