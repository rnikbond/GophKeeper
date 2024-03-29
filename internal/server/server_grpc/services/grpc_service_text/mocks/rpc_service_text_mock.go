// Code generated by MockGen. DO NOT EDIT.
// Source: rpc_service_text.go

// Package grpc_service_text is a generated GoMock package.
package grpc_service_text

import (
	text "GophKeeper/internal/server/model/text"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTextApp is a mock of TextApp interface.
type MockTextApp struct {
	ctrl     *gomock.Controller
	recorder *MockTextAppMockRecorder
}

// MockTextAppMockRecorder is the mock recorder for MockTextApp.
type MockTextAppMockRecorder struct {
	mock *MockTextApp
}

// NewMockTextApp creates a new mock instance.
func NewMockTextApp(ctrl *gomock.Controller) *MockTextApp {
	mock := &MockTextApp{ctrl: ctrl}
	mock.recorder = &MockTextAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTextApp) EXPECT() *MockTextAppMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockTextApp) Change(in text.DataTextFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockTextAppMockRecorder) Change(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockTextApp)(nil).Change), in)
}

// Create mocks base method.
func (m *MockTextApp) Create(in text.DataTextFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTextAppMockRecorder) Create(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTextApp)(nil).Create), in)
}

// Delete mocks base method.
func (m *MockTextApp) Delete(in text.DataTextGet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTextAppMockRecorder) Delete(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTextApp)(nil).Delete), in)
}

// Get mocks base method.
func (m *MockTextApp) Get(in text.DataTextGet) (text.DataTextFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", in)
	ret0, _ := ret[0].(text.DataTextFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockTextAppMockRecorder) Get(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTextApp)(nil).Get), in)
}
