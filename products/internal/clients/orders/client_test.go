package orders

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/NStegura/saga/products/internal/clients/orders/api"
)

// MockOrdersApiClient - mock для OrdersApiClient
type MockOrdersApiClient struct {
	mock.Mock
}

func (m *MockOrdersApiClient) CreateOrder(ctx context.Context, in *api.OrderIn, opts ...grpc.CallOption) (*api.OrderId, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.OrderId), args.Error(1)
}

func (m *MockOrdersApiClient) GetOrderStates(ctx context.Context, in *api.OrderId, opts ...grpc.CallOption) (*api.States, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.States), args.Error(1)
}

func (m *MockOrdersApiClient) GetOrders(ctx context.Context, in *api.UserId, opts ...grpc.CallOption) (*api.Orders, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.Orders), args.Error(1)
}

func (m *MockOrdersApiClient) GetOrder(ctx context.Context, in *api.OrderId, opts ...grpc.CallOption) (*api.OrderOut, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.OrderOut), args.Error(1)
}

func (m *MockOrdersApiClient) GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*api.Pong, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.Pong), args.Error(1)
}

func TestClient_GetProductsToReserve(t *testing.T) {
	mockClient := new(MockOrdersApiClient)

	mockClient.On("GetOrder", mock.Anything, &api.OrderId{OrderId: 123}).
		Return(&api.OrderOut{
			OrderId: 123,
			OrderProducts: []*api.OrderProduct{
				{ProductId: 1, Count: 2},
				{ProductId: 2, Count: 3},
			},
		}, nil)

	client := &Client{
		conn:   nil,
		client: mockClient,
	}

	products, err := client.GetProductsToReserve(123)
	assert.NoError(t, err)
	assert.Len(t, products, 2)
	assert.Equal(t, int64(1), products[0].ProductID)
	assert.Equal(t, int64(2), products[0].Count)
	assert.Equal(t, int64(2), products[1].ProductID)
	assert.Equal(t, int64(3), products[1].Count)

	mockClient.AssertCalled(t, "GetOrder", mock.Anything, &api.OrderId{OrderId: 123})
}
