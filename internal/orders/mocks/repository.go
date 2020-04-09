// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository.go

// Package mock_orders is a generated GoMock package.
package mock_orders

import (
	models "github.com/2020_1_Skycode/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method
func (m *MockRepository) Get(order *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get
func (mr *MockRepositoryMockRecorder) Get(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), order)
}

// InsertOrder mocks base method
func (m *MockRepository) InsertOrder(order *models.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrder", order)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrder indicates an expected call of InsertOrder
func (mr *MockRepositoryMockRecorder) InsertOrder(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrder", reflect.TypeOf((*MockRepository)(nil).InsertOrder), order)
}
