// Code generated by MockGen. DO NOT EDIT.
// Source: ../repository.go

// Package mock_prodtags is a generated GoMock package.
package mock_prodtags

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

// InsertInto mocks base method
func (m *MockRepository) InsertInto(tag *models.ProductTag) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertInto", tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertInto indicates an expected call of InsertInto
func (mr *MockRepositoryMockRecorder) InsertInto(tag interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertInto", reflect.TypeOf((*MockRepository)(nil).InsertInto), tag)
}

// GetByID mocks base method
func (m *MockRepository) GetByID(ID uint64) (*models.ProductTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ID)
	ret0, _ := ret[0].(*models.ProductTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockRepositoryMockRecorder) GetByID(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockRepository)(nil).GetByID), ID)
}

// GetByRestID mocks base method
func (m *MockRepository) GetByRestID(restID uint64) ([]*models.ProductTag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRestID", restID)
	ret0, _ := ret[0].([]*models.ProductTag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRestID indicates an expected call of GetByRestID
func (mr *MockRepositoryMockRecorder) GetByRestID(restID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRestID", reflect.TypeOf((*MockRepository)(nil).GetByRestID), restID)
}

// Delete mocks base method
func (m *MockRepository) Delete(ID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockRepositoryMockRecorder) Delete(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), ID)
}