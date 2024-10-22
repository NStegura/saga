package products

import (
	"context"
	"fmt"

	"github.com/NStegura/saga/tgbot/internal/clients/products/api"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client api.ProductsApiClient

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
		client: api.NewProductsApiClient(conn),
		logger: logger,
	}, nil
}

func (c *Client) GetProducts(ctx context.Context) ([]Product, error) {
	products, err := c.client.GetProducts(ctx, &empty.Empty{})
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %w", err)
	}
	productOut := make([]Product, 0, len(products.Products))
	for _, p := range products.Products {
		productOut = append(productOut, Product{
			ArticleID:   p.ProductId,
			Category:    p.Category,
			Name:        p.Name,
			Description: p.Description,
			Count:       p.Count,
		})
	}
	return productOut, nil
}

func (c *Client) GetProduct(ctx context.Context, articleID int64) (Product, error) {
	product, err := c.client.GetProductInfo(ctx, &api.ProductId{ProductId: articleID})
	if err != nil {
		return Product{}, fmt.Errorf("failed to get product: %w", err)
	}
	return Product{
		ArticleID:   product.ProductId,
		Category:    product.Category,
		Name:        product.Name,
		Description: product.Description,
		Count:       product.Count,
	}, err
}

func (c *Client) Close() {
	_ = c.conn.Close()
}
