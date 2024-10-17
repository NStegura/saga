package payments

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/tgbot/internal/clients/payments/api"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	conn   *grpc.ClientConn
	client api.PaymentsApiClient
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
		client: api.NewPaymentsApiClient(conn),
		logger: logger,
	}, nil
}

func (c *Client) PayOrder(ctx context.Context, orderID int64, status bool) (err error) {
	ctx = metadata.AppendToOutgoingContext(ctx,
		"sender", "tgbot",
		"when", time.Now().Format(time.RFC3339),
	)
	_, err = c.client.UpdatePaymentStatus(ctx, &api.PayStatus{
		OrderId: orderID,
		Status:  status,
	})
	if err != nil {
		return fmt.Errorf("failed to pay order: %w", err)
	}
	return nil
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
