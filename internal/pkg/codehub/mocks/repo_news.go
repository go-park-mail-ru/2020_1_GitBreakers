// Code generated by MockGen. DO NOT EDIT.
// Source: repo_news.go

// Package mock_codehub is a generated GoMock package.
package mock_codehub

import (
	models "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockRepoNewsI is a mock of RepoNewsI interface
type MockRepoNewsI struct {
	ctrl     *gomock.Controller
	recorder *MockRepoNewsIMockRecorder
}

// MockRepoNewsIMockRecorder is the mock recorder for MockRepoNewsI
type MockRepoNewsIMockRecorder struct {
	mock *MockRepoNewsI
}

// NewMockRepoNewsI creates a new mock instance
func NewMockRepoNewsI(ctrl *gomock.Controller) *MockRepoNewsI {
	mock := &MockRepoNewsI{ctrl: ctrl}
	mock.recorder = &MockRepoNewsIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepoNewsI) EXPECT() *MockRepoNewsIMockRecorder {
	return m.recorder
}

// GetNews mocks base method
func (m *MockRepoNewsI) GetNews(repoID, limit, offset int64) (models.NewsSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNews", repoID, limit, offset)
	ret0, _ := ret[0].(models.NewsSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNews indicates an expected call of GetNews
func (mr *MockRepoNewsIMockRecorder) GetNews(repoID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNews", reflect.TypeOf((*MockRepoNewsI)(nil).GetNews), repoID, limit, offset)
}
