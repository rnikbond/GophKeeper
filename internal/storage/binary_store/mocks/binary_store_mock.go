// Code generated by MockGen. DO NOT EDIT.
// Source: binary_store.go

// Package binary_store is a generated GoMock package.
package binary_store

import (
	"GophKeeper/internal/server/model/binary"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBinaryStorage is a mock of BinaryStorage interface.
type MockBinaryStorage struct {
	ctrl     *gomock.Controller
	recorder *MockBinaryStorageMockRecorder
}

// MockBinaryStorageMockRecorder is the mock recorder for MockBinaryStorage.
type MockBinaryStorageMockRecorder struct {
	mock *MockBinaryStorage
}

// NewMockBinaryStorage creates a new mock instance.
func NewMockBinaryStorage(ctrl *gomock.Controller) *MockBinaryStorage {
	mock := &MockBinaryStorage{ctrl: ctrl}
	mock.recorder = &MockBinaryStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBinaryStorage) EXPECT() *MockBinaryStorageMockRecorder {
	return m.recorder
}

// Change mocks base method.
func (m *MockBinaryStorage) Change(in binary.DataFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Change", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Change indicates an expected call of Change.
func (mr *MockBinaryStorageMockRecorder) Change(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Change", reflect.TypeOf((*MockBinaryStorage)(nil).Change), in)
}

// Create mocks base method.
func (m *MockBinaryStorage) Create(in binary.DataFull) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockBinaryStorageMockRecorder) Create(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBinaryStorage)(nil).Create), in)
}

// Delete mocks base method.
func (m *MockBinaryStorage) Delete(in binary.DataGet) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", in)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockBinaryStorageMockRecorder) Delete(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBinaryStorage)(nil).Delete), in)
}

// Get mocks base method.
func (m *MockBinaryStorage) Get(in binary.DataGet) (binary.DataFull, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", in)
	ret0, _ := ret[0].(binary.DataFull)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockBinaryStorageMockRecorder) Get(in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockBinaryStorage)(nil).Get), in)
}
