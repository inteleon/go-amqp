// Code generated by MockGen. DO NOT EDIT.
// Source: amqp/amqp.go

// Package mock_amqp is a generated GoMock package.
package mock_amqp

import (
	gomock "github.com/golang/mock/gomock"
	x "github.com/inteleon/go-amqp/amqp"
	reflect "reflect"
)

// MockAMQP is a mock of AMQP interface
type MockAMQP struct {
	ctrl     *gomock.Controller
	recorder *MockAMQPMockRecorder
}

// MockAMQPMockRecorder is the mock recorder for MockAMQP
type MockAMQPMockRecorder struct {
	mock *MockAMQP
}

// NewMockAMQP creates a new mock instance
func NewMockAMQP(ctrl *gomock.Controller) *MockAMQP {
	mock := &MockAMQP{ctrl: ctrl}
	mock.recorder = &MockAMQPMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAMQP) EXPECT() *MockAMQPMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockAMQP) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockAMQPMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockAMQP)(nil).Connect))
}

// Publish mocks base method
func (m *MockAMQP) Publish(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockAMQPMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockAMQP)(nil).Publish), arg0, arg1)
}

// PublishOnExchange mocks base method
func (m *MockAMQP) PublishOnExchange(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishOnExchange", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishOnExchange indicates an expected call of PublishOnExchange
func (mr *MockAMQPMockRecorder) PublishOnExchange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishOnExchange", reflect.TypeOf((*MockAMQP)(nil).PublishOnExchange), arg0, arg1)
}

// PublishWithHeaders mocks base method
func (m *MockAMQP) PublishWithHeaders(arg0 string, arg1 []byte, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishWithHeaders", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishWithHeaders indicates an expected call of PublishWithHeaders
func (mr *MockAMQPMockRecorder) PublishWithHeaders(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishWithHeaders", reflect.TypeOf((*MockAMQP)(nil).PublishWithHeaders), arg0, arg1, arg2)
}

// PublishOnExchangeWithHeaders mocks base method
func (m *MockAMQP) PublishOnExchangeWithHeaders(arg0 string, arg1 []byte, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishOnExchangeWithHeaders", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishOnExchangeWithHeaders indicates an expected call of PublishOnExchangeWithHeaders
func (mr *MockAMQPMockRecorder) PublishOnExchangeWithHeaders(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishOnExchangeWithHeaders", reflect.TypeOf((*MockAMQP)(nil).PublishOnExchangeWithHeaders), arg0, arg1, arg2)
}

// Consume mocks base method
func (m *MockAMQP) Consume() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume")
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume
func (mr *MockAMQPMockRecorder) Consume() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockAMQP)(nil).Consume))
}

// Ping mocks base method
func (m *MockAMQP) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping
func (mr *MockAMQPMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockAMQP)(nil).Ping))
}

// Close mocks base method
func (m *MockAMQP) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockAMQPMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAMQP)(nil).Close))
}

// Reconnect mocks base method
func (m *MockAMQP) Reconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Reconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Reconnect indicates an expected call of Reconnect
func (mr *MockAMQPMockRecorder) Reconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reconnect", reflect.TypeOf((*MockAMQP)(nil).Reconnect))
}

// MockAMQPClient is a mock of AMQPClient interface
type MockAMQPClient struct {
	ctrl     *gomock.Controller
	recorder *MockAMQPClientMockRecorder
}

// MockAMQPClientMockRecorder is the mock recorder for MockAMQPClient
type MockAMQPClientMockRecorder struct {
	mock *MockAMQPClient
}

// NewMockAMQPClient creates a new mock instance
func NewMockAMQPClient(ctrl *gomock.Controller) *MockAMQPClient {
	mock := &MockAMQPClient{ctrl: ctrl}
	mock.recorder = &MockAMQPClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAMQPClient) EXPECT() *MockAMQPClientMockRecorder {
	return m.recorder
}

// Connect mocks base method
func (m *MockAMQPClient) Connect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect
func (mr *MockAMQPClientMockRecorder) Connect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockAMQPClient)(nil).Connect))
}

// Close mocks base method
func (m *MockAMQPClient) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockAMQPClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockAMQPClient)(nil).Close))
}

// Publish mocks base method
func (m *MockAMQPClient) Publish(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockAMQPClientMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockAMQPClient)(nil).Publish), arg0, arg1)
}

// PublishOnExchange mocks base method
func (m *MockAMQPClient) PublishOnExchange(arg0 string, arg1 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishOnExchange", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishOnExchange indicates an expected call of PublishOnExchange
func (mr *MockAMQPClientMockRecorder) PublishOnExchange(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishOnExchange", reflect.TypeOf((*MockAMQPClient)(nil).PublishOnExchange), arg0, arg1)
}

// PublishWithHeaders mocks base method
func (m *MockAMQPClient) PublishWithHeaders(arg0 string, arg1 []byte, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishWithHeaders", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishWithHeaders indicates an expected call of PublishWithHeaders
func (mr *MockAMQPClientMockRecorder) PublishWithHeaders(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishWithHeaders", reflect.TypeOf((*MockAMQPClient)(nil).PublishWithHeaders), arg0, arg1, arg2)
}

// PublishOnExchangeWithHeaders mocks base method
func (m *MockAMQPClient) PublishOnExchangeWithHeaders(arg0 string, arg1 []byte, arg2 map[string]interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishOnExchangeWithHeaders", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishOnExchangeWithHeaders indicates an expected call of PublishOnExchangeWithHeaders
func (mr *MockAMQPClientMockRecorder) PublishOnExchangeWithHeaders(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishOnExchangeWithHeaders", reflect.TypeOf((*MockAMQPClient)(nil).PublishOnExchangeWithHeaders), arg0, arg1, arg2)
}

// Consume mocks base method
func (m *MockAMQPClient) Consume() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume")
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume
func (mr *MockAMQPClientMockRecorder) Consume() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockAMQPClient)(nil).Consume))
}

// MockAMQPDial is a mock of AMQPDial interface
type MockAMQPDial struct {
	ctrl     *gomock.Controller
	recorder *MockAMQPDialMockRecorder
}

// MockAMQPDialMockRecorder is the mock recorder for MockAMQPDial
type MockAMQPDialMockRecorder struct {
	mock *MockAMQPDial
}

// NewMockAMQPDial creates a new mock instance
func NewMockAMQPDial(ctrl *gomock.Controller) *MockAMQPDial {
	mock := &MockAMQPDial{ctrl: ctrl}
	mock.recorder = &MockAMQPDialMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAMQPDial) EXPECT() *MockAMQPDialMockRecorder {
	return m.recorder
}

// Dial mocks base method
func (m *MockAMQPDial) Dial() (x.AMQPClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Dial")
	ret0, _ := ret[0].(x.AMQPClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Dial indicates an expected call of Dial
func (mr *MockAMQPDialMockRecorder) Dial() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Dial", reflect.TypeOf((*MockAMQPDial)(nil).Dial))
}
