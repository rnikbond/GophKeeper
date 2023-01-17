// Code generated by MockGen. DO NOT EDIT.
// Source: rpc_service_cred.go

// Package grpc_service_cred is a generated GoMock package.
package grpc_service_cred

import (
	cred "GophKeeper/internal/server/model/cred"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCredentialApp is a mock of CredentialApp interface.
type MockCredentialApp struct {
	ctrl     *gomock.Controller
	recorder *MockCredentialAppMockRecorder
}

// MockCredentialAppMockRecorder is the mock recorder for MockCredentialApp.
type MockCredentialAppMockRecorder struct {
	mock *MockCredentialApp
}

// NewMockCredentialApp creates a new mock instance.
func NewMockCredentialApp(ctrl *gomock.Controller) *MockCredentialApp {
	mock := &MockCredentialApp{ctrl: ctrl}
	mock.recorder = &MockCredentialAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredentialApp) EXPECT() *MockCredentialAppMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockCredentialApp) Change(in cred.CredentialFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockCredentialAppMockRecorder) Change(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockCredentialApp)(nil).Change), in)
}

// Create mocks base method.
func (m *MockCredentialApp) Create(in cred.CredentialFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCredentialAppMockRecorder) Create(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCredentialApp)(nil).Create), in)
}

// Delete mocks base method.
func (m *MockCredentialApp) Delete(in cred.CredentialGet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCredentialAppMockRecorder) Delete(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCredentialApp)(nil).Delete), in)
}

// Get mocks base method.
func (m *MockCredentialApp) Get(in cred.CredentialGet) (cred.CredentialFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", in)
	ret0, _ := ret[0].(cred.CredentialFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCredentialAppMockRecorder) Get(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCredentialApp)(nil).Get), in)
}
