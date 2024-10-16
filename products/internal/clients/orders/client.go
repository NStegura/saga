package orders

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/products/internal/clients/orders/orderapi"
	"google.golang.org/grpc/credentials/insecure"

	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn   *grpc.ClientConn
	client orderapi.OrdersApiClient
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
		client: orderapi.NewOrdersApiClient(conn),
	}, nil
}

func (c *Client) GetProductsToReserve(OrderID int64) (ps []Product, err error) {
	ctx := metadata.AppendToOutgoingContext(context.Background(),
		"sender", "productService",
		"when", time.Now().Format(time.RFC3339),
	)
	order, err := c.client.GetOrder(ctx, &orderapi.OrderId{OrderId: OrderID})
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
