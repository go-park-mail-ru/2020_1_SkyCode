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

// GetAllByUserID mocks base method
func (m *MockRepository) GetAllByUserID(userID, count, page uint64) ([]*models.Order, uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", userID, count, page)
	ret0, _ := ret[0].([]*models.Order)
	ret1, _ := ret[1].(uint64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetAllByUserID indicates an expected call of GetAllByUserID
func (mr *MockRepositoryMockRecorder) GetAllByUserID(userID, count, page interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MockRepository)(nil).GetAllByUserID), userID, count, page)
}

// GetByID mocks base method
func (m *MockRepository) GetByID(orderID uint64) (*models.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", orderID)
	ret0, _ := ret[0].(*models.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockRepositoryMockRecorder) GetByID(orderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), orderID)
}

// ChangeStatus mocks base method
func (m *MockRepository) ChangeStatus(orderID uint64, status string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeStatus", orderID, status)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeStatus indicates an expected call of ChangeStatus
func (mr *MockRepositoryMockRecorder) ChangeStatus(orderID, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeStatus", reflect.TypeOf((*MockRepository)(nil).ChangeStatus), orderID, status)
}

// InsertOrder mocks base method
func (m *MockRepository) InsertOrder(order *models.Order, ordProducts []*models.OrderProduct) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertOrder", order, ordProducts)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertOrder indicates an expected call of InsertOrder
func (mr *MockRepositoryMockRecorder) InsertOrder(order, ordProducts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertOrder", reflect.TypeOf((*MockRepository)(nil).InsertOrder), order, ordProducts)
}

// DeleteOrder mocks base method
func (m *MockRepository) DeleteOrder(orderID, userID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteOrder", orderID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteOrder indicates an expected call of DeleteOrder
func (mr *MockRepositoryMockRecorder) DeleteOrder(orderID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteOrder", reflect.TypeOf((*MockRepository)(nil).DeleteOrder), orderID, userID)
}
