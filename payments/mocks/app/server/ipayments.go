// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/server/ipayments.go

// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	reflect "reflect"

	models "github.com/NStegura/saga/payments/internal/services/payment/models"
	gomock "github.com/golang/mock/gomock"
)

// MockPayments is a mock of Payments interface.
type MockPayments struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentsMockRecorder
}

// MockPaymentsMockRecorder is the mock recorder for MockPayments.
type MockPaymentsMockRecorder struct {
	mock *MockPayments
}

// NewMockPayments creates a new mock instance.
func NewMockPayments(ctrl *gomock.Controller) *MockPayments {
	mock := &MockPayments{ctrl: ctrl}
	mock.recorder = &MockPaymentsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPayments) EXPECT() *MockPaymentsMockRecorder {
	return m.recorder
}

// UpdatePaymentStatus mocks base method.
func (m *MockPayments) UpdatePaymentStatus(ctx context.Context, orderID int64, status models.PaymentMessageStatus) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePaymentStatus", ctx, orderID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePaymentStatus indicates an expected call of UpdatePaymentStatus.
func (mr *MockPaymentsMockRecorder) UpdatePaymentStatus(ctx, orderID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePaymentStatus", reflect.TypeOf((*MockPayments)(nil).UpdatePaymentStatus), ctx, orderID, status)
}
