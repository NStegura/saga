package business

import (
	"context"
	"errors"
	"github.com/NStegura/saga/tgbot/internal/clients/orders"
	"github.com/NStegura/saga/tgbot/internal/clients/products"
	"github.com/NStegura/saga/tgbot/internal/domain"
	mock_business "github.com/NStegura/saga/tgbot/mocks/services/business"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBusiness_GetProducts_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	expectedProducts := []products.Product{
		{ArticleID: 1, Category: "Electronics", Name: "Phone", Description: "Smartphone", Count: 10},
		{ArticleID: 2, Category: "Clothing", Name: "T-Shirt", Description: "Cotton T-Shirt", Count: 50},
	}

	mockProductCli.EXPECT().GetProducts(gomock.Any()).Return(expectedProducts, nil)

	result, err := business.GetProducts(context.Background())

	assert.NoError(t, err)
	assert.Len(t, result, len(expectedProducts))
	assert.Equal(t, expectedProducts[0].ArticleID, result[0].ArticleID)
	assert.Equal(t, expectedProducts[1].Category, result[1].Category)
}

func TestBusiness_GetProducts_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	mockProductCli.EXPECT().GetProducts(gomock.Any()).Return(nil, errors.New("internal error"))

	result, err := business.GetProducts(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBusiness_GetProduct_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	expectedProduct := products.Product{
		ArticleID:   1,
		Category:    "Electronics",
		Name:        "Phone",
		Description: "Smartphone",
		Count:       10,
	}

	mockProductCli.EXPECT().GetProduct(gomock.Any(), expectedProduct.ArticleID).Return(expectedProduct, nil)

	result, err := business.GetProduct(context.Background(), expectedProduct.ArticleID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ArticleID, result.ArticleID)
	assert.Equal(t, expectedProduct.Name, result.Name)
}

func TestBusiness_CreateOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderInfo := domain.CreateOrderInfo{
		UserID:      1,
		Description: "New order",
		Products: []domain.OrderProduct{
			{ArticleID: 1, Count: 2},
		},
	}

	mockOrderCli.EXPECT().CreateOrder(gomock.Any(), orderInfo.UserID, orderInfo.Description, gomock.Any()).Return(int64(123), nil)

	orderID, err := business.CreateOrder(context.Background(), orderInfo)

	assert.NoError(t, err)
	assert.Equal(t, int64(123), orderID)
}

func TestBusiness_CreateOrder_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderInfo := domain.CreateOrderInfo{
		UserID:      1,
		Description: "New order",
		Products: []domain.OrderProduct{
			{ArticleID: 1, Count: 2},
		},
	}

	mockOrderCli.EXPECT().CreateOrder(gomock.Any(), orderInfo.UserID, orderInfo.Description, gomock.Any()).Return(int64(0), errors.New("failed to create order"))

	orderID, err := business.CreateOrder(context.Background(), orderInfo)

	assert.Error(t, err)
	assert.Equal(t, int64(0), orderID)
}

func TestBusiness_PayOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderID := int64(123)
	mockPaymentCli.EXPECT().PayOrder(gomock.Any(), orderID, true).Return(nil)

	err := business.PayOrder(context.Background(), orderID, true)

	assert.NoError(t, err)
}

func TestBusiness_PayOrder_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderID := int64(123)
	mockPaymentCli.EXPECT().PayOrder(gomock.Any(), orderID, true).Return(errors.New("payment failed"))

	err := business.PayOrder(context.Background(), orderID, true)

	assert.Error(t, err)
}

func TestBusiness_GetOrder_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderID := int64(123)
	expectedOrder := orders.Order{
		OrderInfo: orders.OrderInfo{
			OrderId:     orderID,
			Description: "Test order",
			State:       "PENDING",
		},
		OrderProducts: []orders.Product{
			{ArticleID: 1, Count: 2},
			{ArticleID: 2, Count: 1},
		},
	}

	expectedProducts := []products.Product{
		{ArticleID: 1, Category: "Electronics", Name: "Phone", Description: "Smartphone", Count: 10},
		{ArticleID: 2, Category: "Clothing", Name: "T-Shirt", Description: "Cotton T-Shirt", Count: 50},
	}

	expectedStatuses := []orders.OrderStatus{
		{Status: "PENDING", Time: time.Now()},
		{Status: "CONFIRMED", Time: time.Now()},
	}

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), orderID).Return(expectedOrder, nil)
	mockProductCli.EXPECT().GetProduct(gomock.Any(), expectedOrder.OrderProducts[0].ArticleID).Return(expectedProducts[0], nil)
	mockProductCli.EXPECT().GetProduct(gomock.Any(), expectedOrder.OrderProducts[1].ArticleID).Return(expectedProducts[1], nil)
	mockOrderCli.EXPECT().GetOrderStatuses(gomock.Any(), orderID).Return(expectedStatuses, nil)

	result, err := business.GetOrder(context.Background(), orderID)

	assert.NoError(t, err)
	assert.Equal(t, orderID, result.OrderID)
	assert.Len(t, result.Products, 2)
	assert.Equal(t, expectedProducts[0].Name, result.Products[0].Name)
	assert.Equal(t, expectedStatuses[0].Status, result.StatusHistory[0].Status)
}

func TestBusiness_GetOrder_Failure_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderID := int64(123)

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), orderID).Return(orders.Order{}, errors.New("failed to get order"))

	result, err := business.GetOrder(context.Background(), orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get orders")
	assert.Equal(t, domain.Order{}, result)
}

func TestBusiness_GetOrder_Failure_GetProduct(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderID := int64(123)
	expectedOrder := orders.Order{
		OrderInfo: orders.OrderInfo{
			OrderId:     orderID,
			Description: "Test order",
			State:       "PENDING",
		},
		OrderProducts: []orders.Product{
			{ArticleID: 1, Count: 2},
			{ArticleID: 2, Count: 1},
		},
	}

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), orderID).Return(expectedOrder, nil)

	mockProductCli.EXPECT().GetProduct(gomock.Any(), expectedOrder.OrderProducts[0].ArticleID).Return(products.Product{}, errors.New("failed to get product"))

	result, err := business.GetOrder(context.Background(), orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get product")
	assert.Equal(t, domain.Order{}, result)
}

func TestBusiness_GetOrder_Failure_GetOrderStatuses(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	orderID := int64(123)
	expectedOrder := orders.Order{
		OrderInfo: orders.OrderInfo{
			OrderId:     orderID,
			Description: "Test order",
			State:       "PENDING",
		},
		OrderProducts: []orders.Product{
			{ArticleID: 1, Count: 2},
			{ArticleID: 2, Count: 1},
		},
	}

	expectedProduct := products.Product{
		ArticleID: 1, Category: "Electronics", Name: "Phone", Description: "Smartphone", Count: 10,
	}

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), orderID).Return(expectedOrder, nil)
	mockProductCli.EXPECT().GetProduct(gomock.Any(), expectedOrder.OrderProducts[0].ArticleID).Return(expectedProduct, nil)
	mockProductCli.EXPECT().GetProduct(gomock.Any(), expectedOrder.OrderProducts[1].ArticleID).Return(expectedProduct, nil)

	mockOrderCli.EXPECT().GetOrderStatuses(gomock.Any(), orderID).Return(nil, errors.New("failed to get statuses"))

	result, err := business.GetOrder(context.Background(), orderID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get order statuses")
	assert.Equal(t, domain.Order{}, result)
}

func TestBusiness_GetOrders_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	userID := int64(456)
	expectedOrderInfos := []orders.OrderInfo{
		{OrderId: 1},
		{OrderId: 2},
	}

	mockOrderCli.EXPECT().GetOrders(gomock.Any(), userID).Return(expectedOrderInfos, nil)

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), expectedOrderInfos[0].OrderId).Return(
		orders.Order{
			OrderInfo: orders.OrderInfo{OrderId: 1},
		}, nil)
	mockOrderCli.EXPECT().GetOrderStatuses(gomock.Any(), expectedOrderInfos[0].OrderId).Return([]orders.OrderStatus{}, nil)
	mockProductCli.EXPECT().GetProduct(gomock.Any(), gomock.Any()).Return(products.Product{}, nil).AnyTimes()

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), expectedOrderInfos[1].OrderId).Return(
		orders.Order{
			OrderInfo: orders.OrderInfo{OrderId: 2},
		}, nil)
	mockOrderCli.EXPECT().GetOrderStatuses(gomock.Any(), expectedOrderInfos[1].OrderId).Return([]orders.OrderStatus{}, nil)

	result, err := business.GetOrders(context.Background(), userID)

	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), result[0].OrderID)
	assert.Equal(t, int64(2), result[1].OrderID)
}

func TestBusiness_GetOrders_Failure_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	userID := int64(456)

	mockOrderCli.EXPECT().GetOrders(gomock.Any(), userID).Return(nil, errors.New("failed to get orders"))

	result, err := business.GetOrders(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get orders")
	assert.Nil(t, result)
}

func TestBusiness_GetOrders_Failure_GetOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProductCli := mock_business.NewMockIProductCli(ctrl)
	mockOrderCli := mock_business.NewMockIOrderCli(ctrl)
	mockPaymentCli := mock_business.NewMockIPaymentCli(ctrl)
	logger := logrus.New()

	business := New(mockOrderCli, mockPaymentCli, mockProductCli, logger)

	userID := int64(456)
	expectedOrderInfos := []orders.OrderInfo{
		{OrderId: 1},
	}

	mockOrderCli.EXPECT().GetOrders(gomock.Any(), userID).Return(expectedOrderInfos, nil)

	mockOrderCli.EXPECT().GetOrder(gomock.Any(), expectedOrderInfos[0].OrderId).Return(orders.Order{}, errors.New("failed to get order"))

	result, err := business.GetOrders(context.Background(), userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get order")
	assert.Nil(t, result)
}
