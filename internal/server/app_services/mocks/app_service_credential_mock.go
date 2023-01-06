// Code generated by MockGen. DO NOT EDIT.
// Source: app_service_credential.go

// Package app_services is a generated GoMock package.
package app_services

import (
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

// Create mocks base method.
func (m *MockCredentialApp) Create(email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCredentialAppMockRecorder) Create(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCredentialApp)(nil).Create), email, password)
}

// Delete mocks base method.
func (m *MockCredentialApp) Delete(email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCredentialAppMockRecorder) Delete(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCredentialApp)(nil).Delete), email)
}

// Read mocks base method.
func (m *MockCredentialApp) Read(email string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", email)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockCredentialAppMockRecorder) Read(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*MockCredentialApp)(nil).Read), email)
}

// Update mocks base method.
func (m *MockCredentialApp) Update(email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCredentialAppMockRecorder) Update(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCredentialApp)(nil).Update), email, password)
}