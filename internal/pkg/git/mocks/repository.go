// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mockGit is a generated GoMock package.
package mockGit

import (
	git "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	permission_types "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockGitRepoI is a mock of GitRepoI interface.
type MockGitRepoI struct {
	ctrl     *gomock.Controller
	recorder *MockGitRepoIMockRecorder
}

// MockGitRepoIMockRecorder is the mock recorder for MockGitRepoI.
type MockGitRepoIMockRecorder struct {
	mock *MockGitRepoI
}

// NewMockGitRepoI creates a new mock instance.
func NewMockGitRepoI(ctrl *gomock.Controller) *MockGitRepoI {
	mock := &MockGitRepoI{ctrl: ctrl}
	mock.recorder = &MockGitRepoIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGitRepoI) EXPECT() *MockGitRepoIMockRecorder {
	return m.recorder
}

// GetByID mocks base method.
func (m *MockGitRepoI) GetByID(id int64) (git.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", id)
	ret0, _ := ret[0].(git.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockGitRepoIMockRecorder) GetByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockGitRepoI)(nil).GetByID), id)
}

// GetByName mocks base method.
func (m *MockGitRepoI) GetByName(userLogin, repoName string) (git.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByName", userLogin, repoName)
	ret0, _ := ret[0].(git.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByName indicates an expected call of GetByName.
func (mr *MockGitRepoIMockRecorder) GetByName(userLogin, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByName", reflect.TypeOf((*MockGitRepoI)(nil).GetByName), userLogin, repoName)
}

// Create mocks base method.
func (m *MockGitRepoI) Create(repos git.Repository) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", repos)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockGitRepoIMockRecorder) Create(repos interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockGitRepoI)(nil).Create), repos)
}

// DeleteByOwnerID mocks base method.
func (m *MockGitRepoI) DeleteByOwnerID(ownerID int64, repoName string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByOwnerID", ownerID, repoName)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByOwnerID indicates an expected call of DeleteByOwnerID.
func (mr *MockGitRepoIMockRecorder) DeleteByOwnerID(ownerID, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByOwnerID", reflect.TypeOf((*MockGitRepoI)(nil).DeleteByOwnerID), ownerID, repoName)
}

// CheckReadAccess mocks base method.
func (m *MockGitRepoI) CheckReadAccess(currentUserId *int64, userLogin, repoName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckReadAccess", currentUserId, userLogin, repoName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckReadAccess indicates an expected call of CheckReadAccess.
func (mr *MockGitRepoIMockRecorder) CheckReadAccess(currentUserId, userLogin, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckReadAccess", reflect.TypeOf((*MockGitRepoI)(nil).CheckReadAccess), currentUserId, userLogin, repoName)
}

// CheckReadAccessById mocks base method.
func (m *MockGitRepoI) CheckReadAccessById(currentUserId *int64, repoId int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckReadAccessById", currentUserId, repoId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckReadAccessById indicates an expected call of CheckReadAccessById.
func (mr *MockGitRepoIMockRecorder) CheckReadAccessById(currentUserId, repoId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckReadAccessById", reflect.TypeOf((*MockGitRepoI)(nil).CheckReadAccessById), currentUserId, repoId)
}

// GetPermission mocks base method.
func (m *MockGitRepoI) GetPermission(currentUserId *int64, userLogin, repoName string) (permission_types.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPermission", currentUserId, userLogin, repoName)
	ret0, _ := ret[0].(permission_types.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPermission indicates an expected call of GetPermission.
func (mr *MockGitRepoIMockRecorder) GetPermission(currentUserId, userLogin, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPermission", reflect.TypeOf((*MockGitRepoI)(nil).GetPermission), currentUserId, userLogin, repoName)
}

// GetPermissionByID mocks base method.
func (m *MockGitRepoI) GetPermissionByID(currentUserId *int64, repoID int64) (permission_types.Permission, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPermissionByID", currentUserId, repoID)
	ret0, _ := ret[0].(permission_types.Permission)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPermissionByID indicates an expected call of GetPermissionByID.
func (mr *MockGitRepoIMockRecorder) GetPermissionByID(currentUserId, repoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPermissionByID", reflect.TypeOf((*MockGitRepoI)(nil).GetPermissionByID), currentUserId, repoID)
}

// GetRepoPathByID mocks base method.
func (m *MockGitRepoI) GetRepoPathByID(repoID int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoPathByID", repoID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepoPathByID indicates an expected call of GetRepoPathByID.
func (mr *MockGitRepoIMockRecorder) GetRepoPathByID(repoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoPathByID", reflect.TypeOf((*MockGitRepoI)(nil).GetRepoPathByID), repoID)
}

// IsRepoExistsByID mocks base method.
func (m *MockGitRepoI) IsRepoExistsByID(repoID int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRepoExistsByID", repoID)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRepoExistsByID indicates an expected call of IsRepoExistsByID.
func (mr *MockGitRepoIMockRecorder) IsRepoExistsByID(repoID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRepoExistsByID", reflect.TypeOf((*MockGitRepoI)(nil).IsRepoExistsByID), repoID)
}

// IsRepoExistsByOwnerId mocks base method.
func (m *MockGitRepoI) IsRepoExistsByOwnerId(ownerId int64, repoName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRepoExistsByOwnerId", ownerId, repoName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRepoExistsByOwnerId indicates an expected call of IsRepoExistsByOwnerId.
func (mr *MockGitRepoIMockRecorder) IsRepoExistsByOwnerId(ownerId, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRepoExistsByOwnerId", reflect.TypeOf((*MockGitRepoI)(nil).IsRepoExistsByOwnerId), ownerId, repoName)
}

// IsRepoExistsByOwnerLogin mocks base method.
func (m *MockGitRepoI) IsRepoExistsByOwnerLogin(ownerLogin, repoName string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsRepoExistsByOwnerLogin", ownerLogin, repoName)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsRepoExistsByOwnerLogin indicates an expected call of IsRepoExistsByOwnerLogin.
func (mr *MockGitRepoIMockRecorder) IsRepoExistsByOwnerLogin(ownerLogin, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsRepoExistsByOwnerLogin", reflect.TypeOf((*MockGitRepoI)(nil).IsRepoExistsByOwnerLogin), ownerLogin, repoName)
}

// GetBranchHashIfExistInRepoByID mocks base method.
func (m *MockGitRepoI) GetBranchHashIfExistInRepoByID(repoID int64, branchName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranchHashIfExistInRepoByID", repoID, branchName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranchHashIfExistInRepoByID indicates an expected call of GetBranchHashIfExistInRepoByID.
func (mr *MockGitRepoIMockRecorder) GetBranchHashIfExistInRepoByID(repoID, branchName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranchHashIfExistInRepoByID", reflect.TypeOf((*MockGitRepoI)(nil).GetBranchHashIfExistInRepoByID), repoID, branchName)
}

// GetBranchInfoByNames mocks base method.
func (m *MockGitRepoI) GetBranchInfoByNames(userLogin, repoName, branchName string) (git.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranchInfoByNames", userLogin, repoName, branchName)
	ret0, _ := ret[0].(git.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranchInfoByNames indicates an expected call of GetBranchInfoByNames.
func (mr *MockGitRepoIMockRecorder) GetBranchInfoByNames(userLogin, repoName, branchName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranchInfoByNames", reflect.TypeOf((*MockGitRepoI)(nil).GetBranchInfoByNames), userLogin, repoName, branchName)
}

// GetBranchesByName mocks base method.
func (m *MockGitRepoI) GetBranchesByName(userLogin, repoName string) ([]git.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBranchesByName", userLogin, repoName)
	ret0, _ := ret[0].([]git.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBranchesByName indicates an expected call of GetBranchesByName.
func (mr *MockGitRepoIMockRecorder) GetBranchesByName(userLogin, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBranchesByName", reflect.TypeOf((*MockGitRepoI)(nil).GetBranchesByName), userLogin, repoName)
}

// GetAnyReposByUserLogin mocks base method.
func (m *MockGitRepoI) GetAnyReposByUserLogin(userLogin string, offset, limit int64) ([]git.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAnyReposByUserLogin", userLogin, offset, limit)
	ret0, _ := ret[0].([]git.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAnyReposByUserLogin indicates an expected call of GetAnyReposByUserLogin.
func (mr *MockGitRepoIMockRecorder) GetAnyReposByUserLogin(userLogin, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAnyReposByUserLogin", reflect.TypeOf((*MockGitRepoI)(nil).GetAnyReposByUserLogin), userLogin, offset, limit)
}

// GetReposByUserLogin mocks base method.
func (m *MockGitRepoI) GetReposByUserLogin(requesterId *int64, userLogin string, offset, limit int64) ([]git.Repository, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReposByUserLogin", requesterId, userLogin, offset, limit)
	ret0, _ := ret[0].([]git.Repository)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReposByUserLogin indicates an expected call of GetReposByUserLogin.
func (mr *MockGitRepoIMockRecorder) GetReposByUserLogin(requesterId, userLogin, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReposByUserLogin", reflect.TypeOf((*MockGitRepoI)(nil).GetReposByUserLogin), requesterId, userLogin, offset, limit)
}

// FilesInCommitByPath mocks base method.
func (m *MockGitRepoI) FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilesInCommitByPath", userLogin, repoName, commitHash, path)
	ret0, _ := ret[0].([]git.FileInCommit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilesInCommitByPath indicates an expected call of FilesInCommitByPath.
func (mr *MockGitRepoIMockRecorder) FilesInCommitByPath(userLogin, repoName, commitHash, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilesInCommitByPath", reflect.TypeOf((*MockGitRepoI)(nil).FilesInCommitByPath), userLogin, repoName, commitHash, path)
}

// GetFileByPath mocks base method.
func (m *MockGitRepoI) GetFileByPath(userLogin, repoName, commitHash, path string) (git.FileCommitted, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileByPath", userLogin, repoName, commitHash, path)
	ret0, _ := ret[0].(git.FileCommitted)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileByPath indicates an expected call of GetFileByPath.
func (mr *MockGitRepoIMockRecorder) GetFileByPath(userLogin, repoName, commitHash, path interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileByPath", reflect.TypeOf((*MockGitRepoI)(nil).GetFileByPath), userLogin, repoName, commitHash, path)
}

// GetCommitsByCommitHash mocks base method.
func (m *MockGitRepoI) GetCommitsByCommitHash(userLogin, repoName, commitHash string, offset, limit int64) ([]git.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitsByCommitHash", userLogin, repoName, commitHash, offset, limit)
	ret0, _ := ret[0].([]git.Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitsByCommitHash indicates an expected call of GetCommitsByCommitHash.
func (mr *MockGitRepoIMockRecorder) GetCommitsByCommitHash(userLogin, repoName, commitHash, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitsByCommitHash", reflect.TypeOf((*MockGitRepoI)(nil).GetCommitsByCommitHash), userLogin, repoName, commitHash, offset, limit)
}

// GetCommitsByBranchName mocks base method.
func (m *MockGitRepoI) GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64) ([]git.Commit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCommitsByBranchName", userLogin, repoName, branchName, offset, limit)
	ret0, _ := ret[0].([]git.Commit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCommitsByBranchName indicates an expected call of GetCommitsByBranchName.
func (mr *MockGitRepoIMockRecorder) GetCommitsByBranchName(userLogin, repoName, branchName, offset, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCommitsByBranchName", reflect.TypeOf((*MockGitRepoI)(nil).GetCommitsByBranchName), userLogin, repoName, branchName, offset, limit)
}

// GetRepoHead mocks base method.
func (m *MockGitRepoI) GetRepoHead(userLogin, repoName string) (git.Branch, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRepoHead", userLogin, repoName)
	ret0, _ := ret[0].(git.Branch)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRepoHead indicates an expected call of GetRepoHead.
func (mr *MockGitRepoIMockRecorder) GetRepoHead(userLogin, repoName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRepoHead", reflect.TypeOf((*MockGitRepoI)(nil).GetRepoHead), userLogin, repoName)
}

// Fork mocks base method.
func (m *MockGitRepoI) Fork(forkRepoName string, userID, repoBaseID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fork", forkRepoName, userID, repoBaseID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Fork indicates an expected call of Fork.
func (mr *MockGitRepoIMockRecorder) Fork(forkRepoName, userID, repoBaseID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fork", reflect.TypeOf((*MockGitRepoI)(nil).Fork), forkRepoName, userID, repoBaseID)
}
