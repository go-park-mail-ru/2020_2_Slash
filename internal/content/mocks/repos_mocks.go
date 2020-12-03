// Code generated by MockGen. DO NOT EDIT.
// Source: internal/content/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "github.com/go-park-mail-ru/2020_2_Slash/internal/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockContentRepository is a mock of ContentRepository interface
type MockContentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockContentRepositoryMockRecorder
}

// MockContentRepositoryMockRecorder is the mock recorder for MockContentRepository
type MockContentRepositoryMockRecorder struct {
	mock *MockContentRepository
}

// NewMockContentRepository creates a new mock instance
func NewMockContentRepository(ctrl *gomock.Controller) *MockContentRepository {
	mock := &MockContentRepository{ctrl: ctrl}
	mock.recorder = &MockContentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContentRepository) EXPECT() *MockContentRepositoryMockRecorder {
	return m.recorder
}

// Insert mocks base method
func (m *MockContentRepository) Insert(content *models.Content) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", content)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert
func (mr *MockContentRepositoryMockRecorder) Insert(content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockContentRepository)(nil).Insert), content)
}

// Update mocks base method
func (m *MockContentRepository) Update(content *models.Content) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", content)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockContentRepositoryMockRecorder) Update(content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockContentRepository)(nil).Update), content)
}

// UpdateImages mocks base method
func (m *MockContentRepository) UpdateImages(content *models.Content) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateImages", content)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateImages indicates an expected call of UpdateImages
func (mr *MockContentRepositoryMockRecorder) UpdateImages(content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateImages", reflect.TypeOf((*MockContentRepository)(nil).UpdateImages), content)
}

// DeleteByID mocks base method
func (m *MockContentRepository) DeleteByID(contentID uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByID", contentID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByID indicates an expected call of DeleteByID
func (mr *MockContentRepositoryMockRecorder) DeleteByID(contentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByID", reflect.TypeOf((*MockContentRepository)(nil).DeleteByID), contentID)
}

// SelectByID mocks base method
func (m *MockContentRepository) SelectByID(contentID uint64) (*models.Content, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectByID", contentID)
	ret0, _ := ret[0].(*models.Content)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectByID indicates an expected call of SelectByID
func (mr *MockContentRepositoryMockRecorder) SelectByID(contentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectByID", reflect.TypeOf((*MockContentRepository)(nil).SelectByID), contentID)
}

// SelectCountriesByID mocks base method
func (m *MockContentRepository) SelectCountriesByID(contentID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectCountriesByID", contentID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectCountriesByID indicates an expected call of SelectCountriesByID
func (mr *MockContentRepositoryMockRecorder) SelectCountriesByID(contentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectCountriesByID", reflect.TypeOf((*MockContentRepository)(nil).SelectCountriesByID), contentID)
}

// SelectGenresByID mocks base method
func (m *MockContentRepository) SelectGenresByID(contentID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectGenresByID", contentID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectGenresByID indicates an expected call of SelectGenresByID
func (mr *MockContentRepositoryMockRecorder) SelectGenresByID(contentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectGenresByID", reflect.TypeOf((*MockContentRepository)(nil).SelectGenresByID), contentID)
}

// SelectActorsByID mocks base method
func (m *MockContentRepository) SelectActorsByID(contentID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectActorsByID", contentID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectActorsByID indicates an expected call of SelectActorsByID
func (mr *MockContentRepositoryMockRecorder) SelectActorsByID(contentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectActorsByID", reflect.TypeOf((*MockContentRepository)(nil).SelectActorsByID), contentID)
}

// SelectDirectorsByID mocks base method
func (m *MockContentRepository) SelectDirectorsByID(contentID uint64) ([]uint64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SelectDirectorsByID", contentID)
	ret0, _ := ret[0].([]uint64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SelectDirectorsByID indicates an expected call of SelectDirectorsByID
func (mr *MockContentRepositoryMockRecorder) SelectDirectorsByID(contentID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectDirectorsByID", reflect.TypeOf((*MockContentRepository)(nil).SelectDirectorsByID), contentID)
}
