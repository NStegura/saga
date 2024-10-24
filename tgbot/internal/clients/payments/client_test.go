package payments

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	"github.com/NStegura/saga/tgbot/internal/clients/payments/api"
)

// MockPaymentsApiClient - mock для PaymentsApiClient
type MockPaymentsApiClient struct {
	mock.Mock
}

func (m *MockPaymentsApiClient) GetPing(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*api.Pong, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*api.Pong), args.Error(1)
}

func (m *MockPaymentsApiClient) UpdatePaymentStatus(ctx context.Context, in *api.PayStatus, opts ...grpc.CallOption) (*empty.Empty, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*empty.Empty), args.Error(1)
}

func TestClient_PayOrder(t *testing.T) {
	mockClient := new(MockPaymentsApiClient)

	orderID := int64(123)
	status := true

	mockClient.On("UpdatePaymentStatus", mock.Anything, &api.PayStatus{
		OrderId: orderID,
		Status:  status,
	}).Return(&empty.Empty{}, nil)

	logger := logrus.New()
	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	err := client.PayOrder(context.Background(), orderID, status)
	assert.NoError(t, err)

	mockClient.AssertCalled(t, "UpdatePaymentStatus", mock.Anything, &api.PayStatus{
		OrderId: orderID,
		Status:  status,
	})
}

func TestClient_PayOrder_False(t *testing.T) {
	mockClient := new(MockPaymentsApiClient)

	orderID := int64(123)
	status := false

	mockClient.On("UpdatePaymentStatus", mock.Anything, &api.PayStatus{
		OrderId: orderID,
		Status:  status,
	}).Return(&empty.Empty{}, nil)

	logger := logrus.New()
	client := &Client{
		conn:   nil,
		client: mockClient,
		logger: logger,
	}

	err := client.PayOrder(context.Background(), orderID, status)
	assert.NoError(t, err)

	mockClient.AssertCalled(t, "UpdatePaymentStatus", mock.Anything, &api.PayStatus{
		OrderId: orderID,
		Status:  status,
	})
}
