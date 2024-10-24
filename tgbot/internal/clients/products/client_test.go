package products

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/NStegura/saga/tgbot/internal/clients/products/api"
)

// MockProductsApiClient - mock для ProductsApiClient
type MockProductsApiClient struct {
	mock.Mock
}

func (m *MockProductsApiClient) GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*api.Pong, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.Pong), args.Error(1)
}

func (m *MockProductsApiClient) GetProducts(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*api.Products, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.Products), args.Error(1)
}

func (m *MockProductsApiClient) GetProductInfo(ctx context.Context, in *api.ProductId, opts ...grpc.CallOption) (*api.Product, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.Product), args.Error(1)
}

func TestClient_GetProducts(t *testing.T) {
	mockClient := new(MockProductsApiClient)
	logger := logrus.New()

	mockClient.On("GetProducts", mock.Anything, &empty.Empty{}).
		Return(&api.Products{
			Products: []*api.Product{
				{ProductId: 1, Category: "Category1", Name: "Product1", Description: "Description1", Count: 10},
				{ProductId: 2, Category: "Category2", Name: "Product2", Description: "Description2", Count: 20},
			},
		}, nil)

	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	products, err := client.GetProducts(context.Background())
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, int64(1), products[0].ArticleID)
	assert.Equal(t, "Category1", products[0].Category)
	assert.Equal(t, "Product1", products[0].Name)
	assert.Equal(t, "Description1", products[0].Description)
	assert.Equal(t, int64(10), products[0].Count)

	mockClient.AssertCalled(t, "GetProducts", mock.Anything, &empty.Empty{})
}

func TestClient_GetProduct(t *testing.T) {
	mockClient := new(MockProductsApiClient)
	logger := logrus.New()

	articleID := int64(1)
	mockClient.On("GetProductInfo", mock.Anything, &api.ProductId{ProductId: articleID}).
		Return(&api.Product{
			ProductId:   articleID,
			Category:    "Category1",
			Name:        "Product1",
			Description: "Description1",
			Count:       10,
		}, nil)

	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	product, err := client.GetProduct(context.Background(), articleID)
	assert.NoError(t, err)
	assert.Equal(t, articleID, product.ArticleID)
	assert.Equal(t, "Category1", product.Category)
	assert.Equal(t, "Product1", product.Name)
	assert.Equal(t, "Description1", product.Description)
	assert.Equal(t, int64(10), product.Count)

	mockClient.AssertCalled(t, "GetProductInfo", mock.Anything, &api.ProductId{ProductId: articleID})
}

func TestClient_GetProduct_Error(t *testing.T) {
	mockClient := new(MockProductsApiClient)
	logger := logrus.New()

	articleID := int64(1)
	mockClient.On("GetProductInfo", mock.Anything, &api.ProductId{ProductId: articleID}).
		Return(&api.Product{}, fmt.Errorf("mock error"))

	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	product, err := client.GetProduct(context.Background(), articleID)
	assert.Error(t, err)
	assert.Equal(t, "failed to get product: mock error", err.Error())
	assert.Equal(t, Product{}, product)

	mockClient.AssertCalled(t, "GetProductInfo", mock.Anything, &api.ProductId{ProductId: articleID})
}
