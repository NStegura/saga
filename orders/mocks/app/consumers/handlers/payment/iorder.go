// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/consumers/handlers/payment/iorder.go

// Package mock_payment is a generated GoMock package.
package mock_payment

import (
	context "context"
	reflect "reflect"

	models "github.com/NStegura/saga/orders/internal/services/order/models"
	gomock "github.com/golang/mock/gomock"
)

// MockOrder is a mock of Order interface.
type MockOrder struct {
	ctrl     *gomock.Controller
	recorder *MockOrderMockRecorder
}

// MockOrderMockRecorder is the mock recorder for MockOrder.
type MockOrderMockRecorder struct {
	mock *MockOrder
}

// NewMockOrder creates a new mock instance.
func NewMockOrder(ctrl *gomock.Controller) *MockOrder {
	mock := &MockOrder{ctrl: ctrl}
	mock.recorder = &MockOrderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrder) EXPECT() *MockOrderMockRecorder {
	return m.recorder
}

// CreateOrderState mocks base method.
func (m *MockOrder) CreateOrderState(ctx context.Context, orderID int64, state models.OrderState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOrderState", ctx, orderID, state)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOrderState indicates an expected call of CreateOrderState.
func (mr *MockOrderMockRecorder) CreateOrderState(ctx, orderID, state interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOrderState", reflect.TypeOf((*MockOrder)(nil).CreateOrderState), ctx, orderID, state)
}