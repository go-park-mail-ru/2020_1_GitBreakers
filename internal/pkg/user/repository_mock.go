// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_user is a generated GoMock package.
package user

import (
	models "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepoUser is a mock of RepoUser interface
type MockRepoUser struct {
	ctrl     *gomock.Controller
	recorder *MockRepoUserMockRecorder
}

// MockRepoUserMockRecorder is the mock recorder for MockRepoUser
type MockRepoUserMockRecorder struct {
	mock *MockRepoUser
}

// NewMockRepoUser creates a new mock instance
func NewMockRepoUser(ctrl *gomock.Controller) *MockRepoUser {
	mock := &MockRepoUser{ctrl: ctrl}
	mock.recorder = &MockRepoUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepoUser) EXPECT() *MockRepoUserMockRecorder {
	return m.recorder
}

// GetUserByIdWithPass mocks base method
func (m *MockRepoUser) GetUserByIdWithPass(id int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByIdWithPass", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByIdWithPass indicates an expected call of GetUserByIdWithPass
func (mr *MockRepoUserMockRecorder) GetUserByIdWithPass(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByIdWithPass", reflect.TypeOf((*MockRepoUser)(nil).GetUserByIdWithPass), id)
}

// GetUserByIdWithoutPass mocks base method
func (m *MockRepoUser) GetUserByIdWithoutPass(id int) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByIdWithoutPass", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByIdWithoutPass indicates an expected call of GetUserByIdWithoutPass
func (mr *MockRepoUserMockRecorder) GetUserByIdWithoutPass(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByIdWithoutPass", reflect.TypeOf((*MockRepoUser)(nil).GetUserByIdWithoutPass), id)
}

// GetUserByLoginWithPass mocks base method
func (m *MockRepoUser) GetUserByLoginWithPass(login string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLoginWithPass", login)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLoginWithPass indicates an expected call of GetUserByLoginWithPass
func (mr *MockRepoUserMockRecorder) GetUserByLoginWithPass(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLoginWithPass", reflect.TypeOf((*MockRepoUser)(nil).GetUserByLoginWithPass), login)
}

// GetByLoginWithoutPass mocks base method
func (m *MockRepoUser) GetByLoginWithoutPass(login string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByLoginWithoutPass", login)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByLoginWithoutPass indicates an expected call of GetByLoginWithoutPass
func (mr *MockRepoUserMockRecorder) GetByLoginWithoutPass(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByLoginWithoutPass", reflect.TypeOf((*MockRepoUser)(nil).GetByLoginWithoutPass), login)
}

// GetLoginByID mocks base method
func (m *MockRepoUser) GetLoginByID(id int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLoginByID", id)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLoginByID indicates an expected call of GetLoginByID
func (mr *MockRepoUserMockRecorder) GetLoginById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLoginByID", reflect.TypeOf((*MockRepoUser)(nil).GetLoginByID), id)
}

// GetIdByLogin mocks base method
func (m *MockRepoUser) GetIdByLogin(login string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIdByLogin", login)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIdByLogin indicates an expected call of GetIdByLogin
func (mr *MockRepoUserMockRecorder) GetIdByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIdByLogin", reflect.TypeOf((*MockRepoUser)(nil).GetIdByLogin), login)
}

// Create mocks base method
func (m *MockRepoUser) Create(newUser models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", newUser)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockRepoUserMockRecorder) Create(newUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepoUser)(nil).Create), newUser)
}

// Update mocks base method
func (m *MockRepoUser) Update(usrUpd models.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", usrUpd)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockRepoUserMockRecorder) Update(usrUpd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepoUser)(nil).Update), usrUpd)
}

// IsExists mocks base method
func (m *MockRepoUser) IsExists(user models.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsExists", user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsExists indicates an expected call of IsExists
func (mr *MockRepoUserMockRecorder) IsExists(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsExists", reflect.TypeOf((*MockRepoUser)(nil).IsExists), user)
}

// DeleteById mocks base method
func (m *MockRepoUser) DeleteById(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteById", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteById indicates an expected call of DeleteById
func (mr *MockRepoUserMockRecorder) DeleteById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteById", reflect.TypeOf((*MockRepoUser)(nil).DeleteById), id)
}

// CheckPass mocks base method
func (m *MockRepoUser) CheckPass(login, newpass string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPass", login, newpass)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPass indicates an expected call of CheckPass
func (mr *MockRepoUserMockRecorder) CheckPass(login, newpass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPass", reflect.TypeOf((*MockRepoUser)(nil).CheckPass), login, newpass)
}

// UploadAvatar mocks base method
func (m *MockRepoUser) UploadAvatar(Name string, Content []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadAvatar", Name, Content)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadAvatar indicates an expected call of UploadAvatar
func (mr *MockRepoUserMockRecorder) UploadAvatar(Name, Content interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadAvatar", reflect.TypeOf((*MockRepoUser)(nil).UploadAvatar), Name, Content)
}

// UpdateAvatarPath mocks base method
func (m *MockRepoUser) UpdateAvatarPath(User models.User, Name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAvatarPath", User, Name)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAvatarPath indicates an expected call of UpdateAvatarPath
func (mr *MockRepoUserMockRecorder) UpdateAvatarPath(User, Name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAvatarPath", reflect.TypeOf((*MockRepoUser)(nil).UpdateAvatarPath), User, Name)
}

// UserCanUpdate mocks base method
func (m *MockRepoUser) UserCanUpdate(user models.User) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserCanUpdate", user)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserCanUpdate indicates an expected call of UserCanUpdate
func (mr *MockRepoUserMockRecorder) UserCanUpdate(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserCanUpdate", reflect.TypeOf((*MockRepoUser)(nil).UserCanUpdate), user)
}
