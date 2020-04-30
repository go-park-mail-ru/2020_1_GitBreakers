// Code generated by MockGen. DO NOT EDIT.
// Source: sess_client.go

// Package mock_interfaces is a generated GoMock package.
package mock_clients

import (
	models "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockSessClientI is a mock of SessClientI interface
type MockSessClientI struct {
	ctrl     *gomock.Controller
	recorder *MockSessClientIMockRecorder
}

// MockSessClientIMockRecorder is the mock recorder for MockSessClientI
type MockSessClientIMockRecorder struct {
	mock *MockSessClientI
}

// NewMockSessClientI creates a new mock instance
func NewMockSessClientI(ctrl *gomock.Controller) *MockSessClientI {
	mock := &MockSessClientI{ctrl: ctrl}
	mock.recorder = &MockSessClientIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessClientI) EXPECT() *MockSessClientIMockRecorder {
	return m.recorder
}

// CreateSess mocks base method
func (m *MockSessClientI) CreateSess(UserID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSess", UserID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSess indicates an expected call of CreateSess
func (mr *MockSessClientIMockRecorder) CreateSess(UserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSess", reflect.TypeOf((*MockSessClientI)(nil).CreateSess), UserID)
}

// DelSess mocks base method
func (m *MockSessClientI) DelSess(SessID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelSess", SessID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelSess indicates an expected call of DelSess
func (mr *MockSessClientIMockRecorder) DelSess(SessID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelSess", reflect.TypeOf((*MockSessClientI)(nil).DelSess), SessID)
}

// GetSess mocks base method
func (m *MockSessClientI) GetSess(SessID string) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSess", SessID)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSess indicates an expected call of GetSess
func (mr *MockSessClientIMockRecorder) GetSess(SessID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSess", reflect.TypeOf((*MockSessClientI)(nil).GetSess), SessID)
}