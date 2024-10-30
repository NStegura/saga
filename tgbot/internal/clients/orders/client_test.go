package orders

import (
	"context"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/NStegura/saga/tgbot/internal/clients/orders/api"
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

func TestClient_GetOrder(t *testing.T) {
	mockClient := new(MockOrdersApiClient)

	mockClient.On("GetOrder", mock.Anything, &api.OrderId{OrderId: 123}).
		Return(&api.OrderOut{
			OrderId:     123,
			Description: "Test Order",
			State:       "Created",
			OrderProducts: []*api.OrderProduct{
				{ProductId: 1, Count: 2},
				{ProductId: 2, Count: 3},
			},
		}, nil)

	logger := logrus.New()
	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	order, err := client.GetOrder(context.Background(), 123)
	assert.NoError(t, err)
	assert.Equal(t, int64(123), order.OrderInfo.OrderId)
	assert.Equal(t, "Test Order", order.OrderInfo.Description)
	assert.Equal(t, "Created", order.OrderInfo.State)
	assert.Len(t, order.OrderProducts, 2)
	assert.Equal(t, int64(1), order.OrderProducts[0].ArticleID)
	assert.Equal(t, int64(2), order.OrderProducts[0].Count)
	assert.Equal(t, int64(2), order.OrderProducts[1].ArticleID)
	assert.Equal(t, int64(3), order.OrderProducts[1].Count)

	mockClient.AssertCalled(t, "GetOrder", mock.Anything, &api.OrderId{OrderId: 123})
}

func TestClient_GetOrders(t *testing.T) {
	mockClient := new(MockOrdersApiClient)

	mockClient.On("GetOrders", mock.Anything, &api.UserId{UserId: 1}).
		Return(&api.Orders{
			Orders: []*api.OrderInfoOut{
				{OrderId: 123, Description: "Test Order 1", State: "Created"},
				{OrderId: 124, Description: "Test Order 2", State: "Completed"},
			},
		}, nil)

	logger := logrus.New()
	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	orders, err := client.GetOrders(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, orders, 2)
	assert.Equal(t, int64(123), orders[0].OrderId)
	assert.Equal(t, "Test Order 1", orders[0].Description)
	assert.Equal(t, "Created", orders[0].State)
	assert.Equal(t, int64(124), orders[1].OrderId)
	assert.Equal(t, "Test Order 2", orders[1].Description)
	assert.Equal(t, "Completed", orders[1].State)

	mockClient.AssertCalled(t, "GetOrders", mock.Anything, &api.UserId{UserId: 1})
}

func TestClient_GetOrderStatuses(t *testing.T) {
	mockClient := new(MockOrdersApiClient)

	mockClient.On("GetOrderStates", mock.Anything, &api.OrderId{OrderId: 123}).
		Return(&api.States{
			OrderStates: []*api.OrderState{
				{State: "Created", Time: timestamppb.New(time.Now())},
				{State: "Confirmed", Time: timestamppb.New(time.Now().Add(60 * time.Second))},
			},
		}, nil)

	logger := logrus.New()
	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	statuses, err := client.GetOrderStatuses(context.Background(), 123)
	assert.NoError(t, err)
	assert.Len(t, statuses, 2)
	assert.Equal(t, "Created", statuses[0].Status)
	assert.NotZero(t, statuses[0].Time)
	assert.Equal(t, "Confirmed", statuses[1].Status)
	assert.NotZero(t, statuses[1].Time)

	mockClient.AssertCalled(t, "GetOrderStates", mock.Anything, &api.OrderId{OrderId: 123})
}

func TestClient_CreateOrder(t *testing.T) {
	mockClient := new(MockOrdersApiClient)

	mockClient.On("CreateOrder", mock.Anything, mock.AnythingOfType("*api.OrderIn")).
		Return(&api.OrderId{OrderId: 123}, nil)

	logger := logrus.New()
	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	orderID, err := client.CreateOrder(context.Background(), 1, "New Order", []Product{
		{ArticleID: 1, Count: 2},
		{ArticleID: 2, Count: 3},
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(123), orderID)

	mockClient.AssertCalled(t, "CreateOrder", mock.Anything, &api.OrderIn{
		UserId:      1,
		Description: "New Order",
		OrderProducts: []*api.OrderProduct{
			{ProductId: 1, Count: 2},
			{ProductId: 2, Count: 3},
		},
	})
}
