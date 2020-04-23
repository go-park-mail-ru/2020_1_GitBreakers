// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock_user is a generated GoMock package.
package mocks

import (
	models "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	"io"
	multipart "mime/multipart"
	"os"
	reflect "reflect"
)

// MockUCUser is a mock of UCUser interface
type MockUCUser struct {
	ctrl     *gomock.Controller
	recorder *MockUCUserMockRecorder
}

// MockUCUserMockRecorder is the mock recorder for MockUCUser
type MockUCUserMockRecorder struct {
	mock *MockUCUser
}

// NewMockUCUser creates a new mock instance
func NewMockUCUser(ctrl *gomock.Controller) *MockUCUser {
	mock := &MockUCUser{ctrl: ctrl}
	mock.recorder = &MockUCUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUCUser) EXPECT() *MockUCUserMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockUCUser) Create(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockUCUserMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUCUser)(nil).Create), user)
}

// Delete mocks base method
func (m *MockUCUser) Delete(user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockUCUserMockRecorder) Delete(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUCUser)(nil).Delete), user)
}

// Update mocks base method
func (m *MockUCUser) Update(userID int64, user models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", userID, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockUCUserMockRecorder) Update(userID, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUCUser)(nil).Update), userID, user)
}

// GetByLogin mocks base method
func (m *MockUCUser) GetByLogin(login string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLogin", login)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLogin indicates an expected call of GetByLogin
func (mr *MockUCUserMockRecorder) GetByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLogin", reflect.TypeOf((*MockUCUser)(nil).GetByLogin), login)
}

// GetByID mocks base method
func (m *MockUCUser) GetByID(userID int64) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", userID)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID
func (mr *MockUCUserMockRecorder) GetByID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockUCUser)(nil).GetByID), userID)
}

// CheckPass mocks base method
func (m *MockUCUser) CheckPass(login, pass string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPass", login, pass)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPass indicates an expected call of CheckPass
func (mr *MockUCUserMockRecorder) CheckPass(login, pass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPass", reflect.TypeOf((*MockUCUser)(nil).CheckPass), login, pass)
}

// UploadAvatar mocks base method
func (m *MockUCUser) UploadAvatar(UserID int64, fileName string, fileData []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadAvatar", UserID, fileName, fileData)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadAvatar indicates an expected call of UploadAvatar
func (mr *MockUCUserMockRecorder) UploadAvatar(User, fileName, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadAvatar", reflect.TypeOf((*MockUCUser)(nil).UploadAvatar), User, fileName, file)
}
