// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_products is a generated GoMock package.
package mock_products

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

// GetRestaurantProducts mocks base method
func (m *MockRepository) GetRestaurantProducts(restID uint64) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRestaurantProducts", restID)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRestaurantProducts indicates an expected call of GetRestaurantProducts
func (mr *MockRepositoryMockRecorder) GetRestaurantProducts(restID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRestaurantProducts", reflect.TypeOf((*MockRepository)(nil).GetRestaurantProducts), restID)
}