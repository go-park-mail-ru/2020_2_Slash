// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_actor is a generated GoMock package.
package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockActorRepository is a mock of ActorRepository interface
type MockActorRepository struct {
	ctrl     *gomock.Controller
	recorder *MockActorRepositoryMockRecorder
}

// MockActorRepositoryMockRecorder is the mock recorder for MockActorRepository
type MockActorRepositoryMockRecorder struct {
	mock *MockActorRepository
}

// NewMockActorRepository creates a new mock instance
func NewMockActorRepository(ctrl *gomock.Controller) *MockActorRepository {
	mock := &MockActorRepository{ctrl: ctrl}
	mock.recorder = &MockActorRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockActorRepository) EXPECT() *MockActorRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockActorRepository) Insert(actor *models.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", actor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockActorRepositoryMockRecorder) Insert(actor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockActorRepository)(nil).Insert), actor)
}

// Update mocks base method
func (m *MockActorRepository) Update(actor *models.Actor) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", actor)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockActorRepositoryMockRecorder) Update(actor interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockActorRepository)(nil).Update), actor)
}

// DeleteById mocks base method
func (m *MockActorRepository) DeleteById(id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById
func (mr *MockActorRepositoryMockRecorder) DeleteById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockActorRepository)(nil).DeleteById), id)
}

// SelectById mocks base method
func (m *MockActorRepository) SelectById(id uint64) (*models.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectById", id)
	ret0, _ := ret[0].(*models.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectById indicates an expected call of SelectById
func (mr *MockActorRepositoryMockRecorder) SelectById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectById", reflect.TypeOf((*MockActorRepository)(nil).SelectById), id)
}

// SelectWhereNameLike mocks base method
func (m *MockActorRepository) SelectWhereNameLike(name string, limit, offset uint64) ([]*models.Actor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectWhereNameLike", name, limit, offset)
	ret0, _ := ret[0].([]*models.Actor)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectWhereNameLike indicates an expected call of SelectWhereNameLike
func (mr *MockActorRepositoryMockRecorder) SelectWhereNameLike(name, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectWhereNameLike", reflect.TypeOf((*MockActorRepository)(nil).SelectWhereNameLike), name, limit, offset)
}
