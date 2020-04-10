// Code generated by MockGen. DO NOT EDIT.
// Source: delivery.go

// Package mock_session is a generated GoMock package.
package mocks

import (
	models "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockSessDelivery is a mock of SessDelivery interface
type MockSessDelivery struct {
	ctrl     *gomock.Controller
	recorder *MockSessDeliveryMockRecorder
}

// MockSessDeliveryMockRecorder is the mock recorder for MockSessDelivery
type MockSessDeliveryMockRecorder struct {
	mock *MockSessDelivery
}

// NewMockSessDelivery creates a new mock instance
func NewMockSessDelivery(ctrl *gomock.Controller) *MockSessDelivery {
	mock := &MockSessDelivery{ctrl: ctrl}
	mock.recorder = &MockSessDeliveryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSessDelivery) EXPECT() *MockSessDeliveryMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockSessDelivery) Create(userID int) (http.Cookie, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", userID)
	ret0, _ := ret[0].(http.Cookie)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create
func (mr *MockSessDeliveryMockRecorder) Create(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessDelivery)(nil).Create), userID)
}

// Delete mocks base method
func (m *MockSessDelivery) Delete(sessID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", sessID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockSessDeliveryMockRecorder) Delete(sessID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessDelivery)(nil).Delete), sessID)
}

// GetBySessID mocks base method
func (m *MockSessDelivery) GetBySessID(sessionID string) (models.Session, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBySessID", sessionID)
	ret0, _ := ret[0].(models.Session)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBySessID indicates an expected call of GetBySessID
func (mr *MockSessDeliveryMockRecorder) GetBySessID(sessionID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBySessID", reflect.TypeOf((*MockSessDelivery)(nil).GetBySessID), sessionID)
}
