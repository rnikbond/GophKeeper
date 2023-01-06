// Code generated by MockGen. DO NOT EDIT.
// Source: auth_store.go

// Package storage is a generated GoMock package.
package storage

import (
	auth "GophKeeper/internal/model/auth"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthStorager is a mock of AuthStorager interface.
type MockAuthStorager struct {
	ctrl     *gomock.Controller
	recorder *MockAuthStoragerMockRecorder
}

// MockAuthStoragerMockRecorder is the mock recorder for MockAuthStorager.
type MockAuthStoragerMockRecorder struct {
	mock *MockAuthStorager
}

// NewMockAuthStorager creates a new mock instance.
func NewMockAuthStorager(ctrl *gomock.Controller) *MockAuthStorager {
	mock := &MockAuthStorager{ctrl: ctrl}
	mock.recorder = &MockAuthStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthStorager) EXPECT() *MockAuthStoragerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockAuthStorager) Create(cred auth.Credential) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cred)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockAuthStoragerMockRecorder) Create(cred interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockAuthStorager)(nil).Create), cred)
}

// Delete mocks base method.
func (m *MockAuthStorager) Delete(email string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", email)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAuthStoragerMockRecorder) Delete(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAuthStorager)(nil).Delete), email)
}

// Find mocks base method.
func (m *MockAuthStorager) Find(email string) (auth.Credential, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Find", email)
	ret0, _ := ret[0].(auth.Credential)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Find indicates an expected call of Find.
func (mr *MockAuthStoragerMockRecorder) Find(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Find", reflect.TypeOf((*MockAuthStorager)(nil).Find), email)
}

// Update mocks base method.
func (m *MockAuthStorager) Update(email, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", email, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAuthStoragerMockRecorder) Update(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAuthStorager)(nil).Update), email, password)
}