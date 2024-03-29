// Code generated by MockGen. DO NOT EDIT.
// Source: credential_store.go

// Package credential_store is a generated GoMock package.
package credential_store

import (
	cred "GophKeeper/internal/server/model/cred"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCredStorage is a mock of CredStorage interface.
type MockCredStorage struct {
	ctrl     *gomock.Controller
	recorder *MockCredStorageMockRecorder
}

// MockCredStorageMockRecorder is the mock recorder for MockCredStorage.
type MockCredStorageMockRecorder struct {
	mock *MockCredStorage
}

// NewMockCredStorage creates a new mock instance.
func NewMockCredStorage(ctrl *gomock.Controller) *MockCredStorage {
	mock := &MockCredStorage{ctrl: ctrl}
	mock.recorder = &MockCredStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCredStorage) EXPECT() *MockCredStorageMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockCredStorage) Change(in cred.CredentialFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockCredStorageMockRecorder) Change(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockCredStorage)(nil).Change), in)
}

// Create mocks base method.
func (m *MockCredStorage) Create(data cred.CredentialFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCredStorageMockRecorder) Create(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCredStorage)(nil).Create), data)
}

// Delete mocks base method.
func (m *MockCredStorage) Delete(in cred.CredentialGet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCredStorageMockRecorder) Delete(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCredStorage)(nil).Delete), in)
}

// Get mocks base method.
func (m *MockCredStorage) Get(in cred.CredentialGet) (cred.CredentialFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", in)
	ret0, _ := ret[0].(cred.CredentialFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCredStorageMockRecorder) Get(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCredStorage)(nil).Get), in)
}
