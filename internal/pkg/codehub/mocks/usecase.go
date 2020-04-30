// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go

// Package mock_codehub is a generated GoMock package.
package mockCodehub

import (
	models "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockUCCodeHubI is a mock of UCCodeHubI interface
type MockUCCodeHubI struct {
	ctrl     *gomock.Controller
	recorder *MockUCCodeHubIMockRecorder
}

// MockUCCodeHubIMockRecorder is the mock recorder for MockUCCodeHubI
type MockUCCodeHubIMockRecorder struct {
	mock *MockUCCodeHubI
}

// NewMockUCCodeHubI creates a new mock instance
func NewMockUCCodeHubI(ctrl *gomock.Controller) *MockUCCodeHubI {
	mock := &MockUCCodeHubI{ctrl: ctrl}
	mock.recorder = &MockUCCodeHubIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUCCodeHubI) EXPECT() *MockUCCodeHubIMockRecorder {
	return m.recorder
}

// ModifyStar mocks base method
func (m *MockUCCodeHubI) ModifyStar(star models.Star) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ModifyStar", star)
	ret0, _ := ret[0].(error)
	return ret0
}

// ModifyStar indicates an expected call of ModifyStar
func (mr *MockUCCodeHubIMockRecorder) ModifyStar(star interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ModifyStar", reflect.TypeOf((*MockUCCodeHubI)(nil).ModifyStar), star)
}

// GetStarredRepos mocks base method
func (m *MockUCCodeHubI) GetStarredRepos(userID, limit, offset int64) (models.RepoSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStarredRepos", userID, limit, offset)
	ret0, _ := ret[0].(models.RepoSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStarredRepos indicates an expected call of GetStarredRepos
func (mr *MockUCCodeHubIMockRecorder) GetStarredRepos(userID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStarredRepos", reflect.TypeOf((*MockUCCodeHubI)(nil).GetStarredRepos), userID, limit, offset)
}

// CreateIssue mocks base method
func (m *MockUCCodeHubI) CreateIssue(issue models.Issue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateIssue", issue)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateIssue indicates an expected call of CreateIssue
func (mr *MockUCCodeHubIMockRecorder) CreateIssue(issue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateIssue", reflect.TypeOf((*MockUCCodeHubI)(nil).CreateIssue), issue)
}

// UpdateIssue mocks base method
func (m *MockUCCodeHubI) UpdateIssue(issue models.Issue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateIssue", issue)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateIssue indicates an expected call of UpdateIssue
func (mr *MockUCCodeHubIMockRecorder) UpdateIssue(issue interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateIssue", reflect.TypeOf((*MockUCCodeHubI)(nil).UpdateIssue), issue)
}

// CloseIssue mocks base method
func (m *MockUCCodeHubI) CloseIssue(issueID, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseIssue", issueID, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseIssue indicates an expected call of CloseIssue
func (mr *MockUCCodeHubIMockRecorder) CloseIssue(issueID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseIssue", reflect.TypeOf((*MockUCCodeHubI)(nil).CloseIssue), issueID, userID)
}

// GetIssuesList mocks base method
func (m *MockUCCodeHubI) GetIssuesList(repoID, userID, limit, offset int64) (models.IssuesSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIssuesList", repoID, userID, limit, offset)
	ret0, _ := ret[0].(models.IssuesSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIssuesList indicates an expected call of GetIssuesList
func (mr *MockUCCodeHubIMockRecorder) GetIssuesList(repoID, userID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIssuesList", reflect.TypeOf((*MockUCCodeHubI)(nil).GetIssuesList), repoID, userID, limit, offset)
}

// GetIssue mocks base method
func (m *MockUCCodeHubI) GetIssue(issueID, userID int64) (models.Issue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIssue", issueID, userID)
	ret0, _ := ret[0].(models.Issue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIssue indicates an expected call of GetIssue
func (mr *MockUCCodeHubIMockRecorder) GetIssue(issueID, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIssue", reflect.TypeOf((*MockUCCodeHubI)(nil).GetIssue), issueID, userID)
}

// GetNews mocks base method
func (m *MockUCCodeHubI) GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNews", repoID, userID, limit, offset)
	ret0, _ := ret[0].(models.NewsSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNews indicates an expected call of GetNews
func (mr *MockUCCodeHubIMockRecorder) GetNews(repoID, userID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNews", reflect.TypeOf((*MockUCCodeHubI)(nil).GetNews), repoID, userID, limit, offset)
}

// GetUserStaredList mocks base method
func (m *MockUCCodeHubI) GetUserStaredList(repoID, limit, offset int64) (models.UserSet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserStaredList", repoID, limit, offset)
	ret0, _ := ret[0].(models.UserSet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserStaredList indicates an expected call of GetUserStaredList
func (mr *MockUCCodeHubIMockRecorder) GetUserStaredList(repoID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserStaredList", reflect.TypeOf((*MockUCCodeHubI)(nil).GetUserStaredList), repoID, limit, offset)
}
