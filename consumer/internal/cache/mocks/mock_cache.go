// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "consumer/internal/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCache is a mock of Cache interface.
type MockCache struct {
	ctrl     *gomock.Controller
	recorder *MockCacheMockRecorder
}

// MockCacheMockRecorder is the mock recorder for MockCache.
type MockCacheMockRecorder struct {
	mock *MockCache
}

// NewMockCache creates a new mock instance.
func NewMockCache(ctrl *gomock.Controller) *MockCache {
	mock := &MockCache{ctrl: ctrl}
	mock.recorder = &MockCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCache) EXPECT() *MockCacheMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockCache) Add(order models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockCacheMockRecorder) Add(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockCache)(nil).Add), order)
}

// GC mocks base method.
func (m *MockCache) GC() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "GC")
}

// GC indicates an expected call of GC.
func (mr *MockCacheMockRecorder) GC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GC", reflect.TypeOf((*MockCache)(nil).GC))
}

// Get mocks base method.
func (m *MockCache) Get(OrderUID string) (models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", OrderUID)
	ret0, _ := ret[0].(models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCacheMockRecorder) Get(OrderUID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCache)(nil).Get), OrderUID)
}

// GetAll mocks base method.
func (m *MockCache) GetAll() []models.Order {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]models.Order)
	return ret0
}

// GetAll indicates an expected call of GetAll.
func (mr *MockCacheMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockCache)(nil).GetAll))
}

// StartGC mocks base method.
func (m *MockCache) StartGC() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StartGC")
}

// StartGC indicates an expected call of StartGC.
func (mr *MockCacheMockRecorder) StartGC() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartGC", reflect.TypeOf((*MockCache)(nil).StartGC))
}

// deleteExpiredKeys mocks base method.
func (m *MockCache) DeleteExpiredKeys() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "deleteExpiredKeys")
}

// deleteExpiredKeys indicates an expected call of deleteExpiredKeys.
func (mr *MockCacheMockRecorder) DeleteExpiredKeys() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "deleteExpiredKeys", reflect.TypeOf((*MockCache)(nil).DeleteExpiredKeys))
}
