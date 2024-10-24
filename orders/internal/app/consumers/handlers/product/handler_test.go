package product

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/IBM/sarama"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/NStegura/saga/orders/internal/app/consumers/handlers/product/models"
	"github.com/NStegura/saga/orders/internal/clients/redis"
	orderModels "github.com/NStegura/saga/orders/internal/services/order/models"
	mock_product "github.com/NStegura/saga/orders/mocks/app/consumers/handlers/product"
)

type MockConsumerGroupSession struct {
	mock.Mock
}

func (m *MockConsumerGroupSession) Claims() map[string][]int32 {
	panic("implement me")
}

func (m *MockConsumerGroupSession) MemberID() string {
	panic("implement me")
}

func (m *MockConsumerGroupSession) GenerationID() int32 {
	panic("implement me")
}

func (m *MockConsumerGroupSession) MarkOffset(topic string, partition int32, offset int64, metadata string) {
	panic("implement me")
}

func (m *MockConsumerGroupSession) Commit() {
	panic("implement me")
}

func (m *MockConsumerGroupSession) ResetOffset(topic string, partition int32, offset int64, metadata string) {
	panic("implement me")
}

func (m *MockConsumerGroupSession) MarkMessage(msg *sarama.ConsumerMessage, metadata string) {
	m.Called(msg, metadata)
}

func (m *MockConsumerGroupSession) Context() context.Context {
	args := m.Called()
	return args.Get(0).(context.Context)
}

type MockConsumerGroupClaim struct {
	mock.Mock
}

func (m *MockConsumerGroupClaim) Topic() string {
	panic("implement me")
}

func (m *MockConsumerGroupClaim) Partition() int32 {
	panic("implement me")
}

func (m *MockConsumerGroupClaim) InitialOffset() int64 {
	panic("implement me")
}

func (m *MockConsumerGroupClaim) HighWaterMarkOffset() int64 {
	panic("implement me")
}

func (m *MockConsumerGroupClaim) Messages() <-chan *sarama.ConsumerMessage {
	args := m.Called()
	return args.Get(0).(<-chan *sarama.ConsumerMessage)
}

func mockMessageChan(msgBytes []byte) <-chan *sarama.ConsumerMessage {
	msgChan := make(chan *sarama.ConsumerMessage, 1)
	msgChan <- &sarama.ConsumerMessage{
		Value: msgBytes,
	}
	close(msgChan)
	return msgChan
}

func TestIncomeHandler_ConsumeClaim_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrder := mock_product.NewMockOrder(ctrl)
	mockCache := mock_product.NewMockCache(ctrl)
	logger := logrus.New()

	handler := New(mockOrder, mockCache, logger)

	mockSession := new(MockConsumerGroupSession)
	mockClaim := new(MockConsumerGroupClaim)

	key := uuid.New()
	message := models.ProductMessage{
		OrderID: 12345,
		Status:  models.CREATED,
		IKey:    key,
	}

	msgBytes, _ := json.Marshal(message)

	mockClaim.On("Messages").Return(mockMessageChan(msgBytes)).Once()
	mockSession.On("MarkMessage", mock.Anything, "").Return().Once()

	mockCache.EXPECT().Get(context.Background(), message.IKey).Return(redis.ErrCacheMiss).Times(1)
	mockOrder.EXPECT().CreateOrderState(context.Background(), message.OrderID, orderModels.RESERVE_CREATED).Return(nil).Times(1)
	mockCache.EXPECT().Set(context.Background(), message.IKey).Return(nil).Times(1)

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)

	mockSession.AssertExpectations(t)
	mockClaim.AssertExpectations(t)
}

func TestIncomeHandler_ConsumeClaim_IdempotentKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrder := mock_product.NewMockOrder(ctrl)
	mockCache := mock_product.NewMockCache(ctrl)
	logger := logrus.New()

	handler := New(mockOrder, mockCache, logger)

	mockSession := new(MockConsumerGroupSession)
	mockClaim := new(MockConsumerGroupClaim)

	key := uuid.New()
	message := models.ProductMessage{
		OrderID: 12345,
		Status:  models.CREATED,
		IKey:    key,
	}

	msgBytes, _ := json.Marshal(message)

	mockClaim.On("Messages").Return(mockMessageChan(msgBytes)).Once()
	mockSession.On("MarkMessage", mock.Anything, "").Return().Once()

	mockCache.EXPECT().Get(context.Background(), message.IKey).Return(nil).Times(1)

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)
	mockSession.AssertNotCalled(t, "MarkMessage")
}

func TestIncomeHandler_ConsumeClaim_UnknownStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrder := mock_product.NewMockOrder(ctrl)
	mockCache := mock_product.NewMockCache(ctrl)
	logger := logrus.New()

	handler := New(mockOrder, mockCache, logger)

	mockSession := new(MockConsumerGroupSession)
	mockClaim := new(MockConsumerGroupClaim)

	key := uuid.New()
	message := models.ProductMessage{
		OrderID: 12345,
		Status:  "unknown-status",
		IKey:    key,
	}

	msgBytes, _ := json.Marshal(message)

	mockClaim.On("Messages").Return(mockMessageChan(msgBytes)).Once()
	mockSession.On("MarkMessage", mock.Anything, "").Return().Once()

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)
	mockSession.AssertNotCalled(t, "MarkMessage")
}

func TestIncomeHandler_ConsumeClaim_CreateOrderStateError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockOrder := mock_product.NewMockOrder(ctrl)
	mockCache := mock_product.NewMockCache(ctrl)
	logger := logrus.New()

	handler := New(mockOrder, mockCache, logger)

	mockSession := new(MockConsumerGroupSession)
	mockClaim := new(MockConsumerGroupClaim)

	key := uuid.New()
	message := models.ProductMessage{
		OrderID: 12345,
		Status:  models.CREATED,
		IKey:    key,
	}

	msgBytes, _ := json.Marshal(message)

	mockClaim.On("Messages").Return(mockMessageChan(msgBytes)).Once()

	mockCache.EXPECT().Get(context.Background(), message.IKey).Return(redis.ErrCacheMiss).Times(1)
	mockOrder.EXPECT().CreateOrderState(
		context.Background(), message.OrderID, orderModels.RESERVE_CREATED,
	).Return(errors.New("create order state error")).Times(1)

	err := handler.ConsumeClaim(mockSession, mockClaim)
	assert.NoError(t, err)

	mockSession.AssertNotCalled(t, "MarkMessage")
}
