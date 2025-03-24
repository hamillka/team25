// Code generated by MockGen. DO NOT EDIT.
// Source: office.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/hamillka/team25/backend/internal/models"
)

// MockOfficeRepository is a mock of OfficeRepository interface.
type MockOfficeRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOfficeRepositoryMockRecorder
}

// MockOfficeRepositoryMockRecorder is the mock recorder for MockOfficeRepository.
type MockOfficeRepositoryMockRecorder struct {
	mock *MockOfficeRepository
}

// NewMockOfficeRepository creates a new mock instance.
func NewMockOfficeRepository(ctrl *gomock.Controller) *MockOfficeRepository {
	mock := &MockOfficeRepository{ctrl: ctrl}
	mock.recorder = &MockOfficeRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOfficeRepository) EXPECT() *MockOfficeRepositoryMockRecorder {
	return m.recorder
}

// AddOffice mocks base method.
func (m *MockOfficeRepository) AddOffice(number, floor int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOffice", number, floor)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOffice indicates an expected call of AddOffice.
func (mr *MockOfficeRepositoryMockRecorder) AddOffice(number, floor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOffice", reflect.TypeOf((*MockOfficeRepository)(nil).AddOffice), number, floor)
}

// EditOffice mocks base method.
func (m *MockOfficeRepository) EditOffice(id, number, floor int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditOffice", id, number, floor)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EditOffice indicates an expected call of EditOffice.
func (mr *MockOfficeRepositoryMockRecorder) EditOffice(id, number, floor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditOffice", reflect.TypeOf((*MockOfficeRepository)(nil).EditOffice), id, number, floor)
}

// GetAllOffices mocks base method.
func (m *MockOfficeRepository) GetAllOffices() ([]models.Office, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllOffices")
	ret0, _ := ret[0].([]models.Office)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllOffices indicates an expected call of GetAllOffices.
func (mr *MockOfficeRepositoryMockRecorder) GetAllOffices() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllOffices", reflect.TypeOf((*MockOfficeRepository)(nil).GetAllOffices))
}

// GetOfficeByID mocks base method.
func (m *MockOfficeRepository) GetOfficeByID(id int64) (models.Office, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOfficeByID", id)
	ret0, _ := ret[0].(models.Office)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOfficeByID indicates an expected call of GetOfficeByID.
func (mr *MockOfficeRepositoryMockRecorder) GetOfficeByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOfficeByID", reflect.TypeOf((*MockOfficeRepository)(nil).GetOfficeByID), id)
}
