// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/hamillka/team25/backend/internal/models"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// CheckUserRole mocks base method.
func (m *MockUserRepository) CheckUserRole(id int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserRole", id)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserRole indicates an expected call of CheckUserRole.
func (mr *MockUserRepositoryMockRecorder) CheckUserRole(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserRole", reflect.TypeOf((*MockUserRepository)(nil).CheckUserRole), id)
}

// CreateUser mocks base method.
func (m *MockUserRepository) CreateUser(login, password string, role int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", login, password, role)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepositoryMockRecorder) CreateUser(login, password, role interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepository)(nil).CreateUser), login, password, role)
}

// CreateUserDoctor mocks base method.
func (m *MockUserRepository) CreateUserDoctor(login, password string, role, docID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserDoctor", login, password, role, docID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserDoctor indicates an expected call of CreateUserDoctor.
func (mr *MockUserRepositoryMockRecorder) CreateUserDoctor(login, password, role, docID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserDoctor", reflect.TypeOf((*MockUserRepository)(nil).CreateUserDoctor), login, password, role, docID)
}

// CreateUserPatient mocks base method.
func (m *MockUserRepository) CreateUserPatient(login, password string, role, patID int64) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUserPatient", login, password, role, patID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUserPatient indicates an expected call of CreateUserPatient.
func (mr *MockUserRepositoryMockRecorder) CreateUserPatient(login, password, role, patID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUserPatient", reflect.TypeOf((*MockUserRepository)(nil).CreateUserPatient), login, password, role, patID)
}

// GetUserByLogin mocks base method.
func (m *MockUserRepository) GetUserByLogin(login string) (*models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", login)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockUserRepositoryMockRecorder) GetUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUserRepository)(nil).GetUserByLogin), login)
}

// GetUserByLoginAndPassword mocks base method.
func (m *MockUserRepository) GetUserByLoginAndPassword(login, password string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLoginAndPassword", login, password)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLoginAndPassword indicates an expected call of GetUserByLoginAndPassword.
func (mr *MockUserRepositoryMockRecorder) GetUserByLoginAndPassword(login, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLoginAndPassword", reflect.TypeOf((*MockUserRepository)(nil).GetUserByLoginAndPassword), login, password)
}
