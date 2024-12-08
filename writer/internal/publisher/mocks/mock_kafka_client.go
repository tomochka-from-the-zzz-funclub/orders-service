// Code generated by MockGen. DO NOT EDIT.
// Source: internal/publisher/interface.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	models "writer/internal/models"

	gomock "github.com/golang/mock/gomock"
)

// MockInterfaceKafkaClient is a mock of InterfaceKafkaClient interface.
type MockInterfaceKafkaClient struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceKafkaClientMockRecorder
}

// MockInterfaceKafkaClientMockRecorder is the mock recorder for MockInterfaceKafkaClient.
type MockInterfaceKafkaClientMockRecorder struct {
	mock *MockInterfaceKafkaClient
}

// NewMockInterfaceKafkaClient creates a new mock instance.
func NewMockInterfaceKafkaClient(ctrl *gomock.Controller) *MockInterfaceKafkaClient {
	mock := &MockInterfaceKafkaClient{ctrl: ctrl}
	mock.recorder = &MockInterfaceKafkaClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterfaceKafkaClient) EXPECT() *MockInterfaceKafkaClientMockRecorder {
	return m.recorder
}

// SendOrderToKafka mocks base method.
func (m *MockInterfaceKafkaClient) SendOrderToKafka(topic string, message models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendOrderToKafka", topic, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendOrderToKafka indicates an expected call of SendOrderToKafka.
func (mr *MockInterfaceKafkaClientMockRecorder) SendOrderToKafka(topic, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendOrderToKafka", reflect.TypeOf((*MockInterfaceKafkaClient)(nil).SendOrderToKafka), topic, message)
}
