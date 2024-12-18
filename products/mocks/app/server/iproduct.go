// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/server/iproduct.go

// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	reflect "reflect"

	models "github.com/NStegura/saga/products/internal/services/product/models"
	gomock "github.com/golang/mock/gomock"
)

// MockProduct is a mock of Product interface.
type MockProduct struct {
	ctrl     *gomock.Controller
	recorder *MockProductMockRecorder
}

// MockProductMockRecorder is the mock recorder for MockProduct.
type MockProductMockRecorder struct {
	mock *MockProduct
}

// NewMockProduct creates a new mock instance.
func NewMockProduct(ctrl *gomock.Controller) *MockProduct {
	mock := &MockProduct{ctrl: ctrl}
	mock.recorder = &MockProductMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProduct) EXPECT() *MockProductMockRecorder {
	return m.recorder
}

// GetProductInfo mocks base method.
func (m *MockProduct) GetProductInfo(ctx context.Context, productID int64) (models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductInfo", ctx, productID)
	ret0, _ := ret[0].(models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductInfo indicates an expected call of GetProductInfo.
func (mr *MockProductMockRecorder) GetProductInfo(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductInfo", reflect.TypeOf((*MockProduct)(nil).GetProductInfo), ctx, productID)
}

// GetProducts mocks base method.
func (m *MockProduct) GetProducts(arg0 context.Context) ([]models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", arg0)
	ret0, _ := ret[0].([]models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts.
func (mr *MockProductMockRecorder) GetProducts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockProduct)(nil).GetProducts), arg0)
}
