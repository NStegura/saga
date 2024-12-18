// Code generated by MockGen. DO NOT EDIT.
// Source: internal/app/server/isystem.go

// Package mock_server is a generated GoMock package.
package mock_server

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSystem is a mock of System interface.
type MockSystem struct {
	ctrl     *gomock.Controller
	recorder *MockSystemMockRecorder
}

// MockSystemMockRecorder is the mock recorder for MockSystem.
type MockSystemMockRecorder struct {
	mock *MockSystem
}

// NewMockSystem creates a new mock instance.
func NewMockSystem(ctrl *gomock.Controller) *MockSystem {
	mock := &MockSystem{ctrl: ctrl}
	mock.recorder = &MockSystemMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSystem) EXPECT() *MockSystemMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockSystem) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockSystemMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockSystem)(nil).Ping), arg0)
}
