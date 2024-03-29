// Code generated by MockGen. DO NOT EDIT.
// Source: rpc_service_card.go

// Package grpc_service_card is a generated GoMock package.
package grpc_service_card

import (
	card "GophKeeper/internal/server/model/card"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCardApp is a mock of CardApp interface.
type MockCardApp struct {
	ctrl     *gomock.Controller
	recorder *MockCardAppMockRecorder
}

// MockCardAppMockRecorder is the mock recorder for MockCardApp.
type MockCardAppMockRecorder struct {
	mock *MockCardApp
}

// NewMockCardApp creates a new mock instance.
func NewMockCardApp(ctrl *gomock.Controller) *MockCardApp {
	mock := &MockCardApp{ctrl: ctrl}
	mock.recorder = &MockCardAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCardApp) EXPECT() *MockCardAppMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockCardApp) Change(in card.DataCardFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockCardAppMockRecorder) Change(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockCardApp)(nil).Change), in)
}

// Create mocks base method.
func (m *MockCardApp) Create(data card.DataCardFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCardAppMockRecorder) Create(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCardApp)(nil).Create), data)
}

// Delete mocks base method.
func (m *MockCardApp) Delete(in card.DataCardGet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCardAppMockRecorder) Delete(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCardApp)(nil).Delete), in)
}

// Get mocks base method.
func (m *MockCardApp) Get(in card.DataCardGet) (card.DataCardFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", in)
	ret0, _ := ret[0].(card.DataCardFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCardAppMockRecorder) Get(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCardApp)(nil).Get), in)
}
