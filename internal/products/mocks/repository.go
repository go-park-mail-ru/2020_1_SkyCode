// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository.go

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

// GetProductsByRestID mocks base method
func (m *MockRepository) GetProductsByRestID(restID uint64) ([]*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsByRestID", restID)
	ret0, _ := ret[0].([]*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsByRestID indicates an expected call of GetProductsByRestID
func (mr *MockRepositoryMockRecorder) GetProductsByRestID(restID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsByRestID", reflect.TypeOf((*MockRepository)(nil).GetProductsByRestID), restID)
}

// GetProductByID mocks base method
func (m *MockRepository) GetProductByID(prodID uint64) (*models.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductByID", prodID)
	ret0, _ := ret[0].(*models.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductByID indicates an expected call of GetProductByID
func (mr *MockRepositoryMockRecorder) GetProductByID(prodID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductByID", reflect.TypeOf((*MockRepository)(nil).GetProductByID), prodID)
}

// InsertInto mocks base method
func (m *MockRepository) InsertInto(product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInto", product)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInto indicates an expected call of InsertInto
func (mr *MockRepositoryMockRecorder) InsertInto(product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInto", reflect.TypeOf((*MockRepository)(nil).InsertInto), product)
}

// Update mocks base method
func (m *MockRepository) Update(product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", product)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockRepositoryMockRecorder) Update(product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), product)
}

// UpdateImage mocks base method
func (m *MockRepository) UpdateImage(product *models.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateImage", product)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateImage indicates an expected call of UpdateImage
func (mr *MockRepositoryMockRecorder) UpdateImage(product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateImage", reflect.TypeOf((*MockRepository)(nil).UpdateImage), product)
}

// Delete mocks base method
func (m *MockRepository) Delete(prodID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", prodID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(prodID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), prodID)
}
