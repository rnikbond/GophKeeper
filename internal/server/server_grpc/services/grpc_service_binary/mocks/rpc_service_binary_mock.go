// Code generated by MockGen. DO NOT EDIT.
// Source: rpc_service_binary.go

// Package grpc_service_binary is a generated GoMock package.
package grpc_service_binary

import (
	"GophKeeper/internal/server/model/binary"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBinaryApp is a mock of BinaryApp interface.
type MockBinaryApp struct {
	ctrl     *gomock.Controller
	recorder *MockBinaryAppMockRecorder
}

// MockBinaryAppMockRecorder is the mock recorder for MockBinaryApp.
type MockBinaryAppMockRecorder struct {
	mock *MockBinaryApp
}

// NewMockBinaryApp creates a new mock instance.
func NewMockBinaryApp(ctrl *gomock.Controller) *MockBinaryApp {
	mock := &MockBinaryApp{ctrl: ctrl}
	mock.recorder = &MockBinaryAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBinaryApp) EXPECT() *MockBinaryAppMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockBinaryApp) Change(in binary.DataFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockBinaryAppMockRecorder) Change(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockBinaryApp)(nil).Change), in)
}

// Create mocks base method.
func (m *MockBinaryApp) Create(in binary.DataFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockBinaryAppMockRecorder) Create(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBinaryApp)(nil).Create), in)
}

// Delete mocks base method.
func (m *MockBinaryApp) Delete(in binary.DataGet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBinaryAppMockRecorder) Delete(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBinaryApp)(nil).Delete), in)
}

// Get mocks base method.
func (m *MockBinaryApp) Get(in binary.DataGet) (binary.DataFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", in)
	ret0, _ := ret[0].(binary.DataFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBinaryAppMockRecorder) Get(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBinaryApp)(nil).Get), in)
}
