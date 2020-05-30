package git

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type GitUseCaseI interface {
	Create(userID int64, repos *gitmodels.Repository) error
	GetRepo(userName string, repoName string, requestUserID *int64) (gitmodels.Repository, error)
	DeleteByOwnerID(ownerID int64, repoName string) error
	GetRepoList(userName string, offset, limit int64, requestUserID *int64) (gitmodels.RepositorySet, error)
	GetBranchList(requestUserID *int64, userName string, repoName string) (gitmodels.BranchSet, error)
	FilesInCommitByPath(requset gitmodels.FilesCommitRequest, requesrUserID *int64) (gitmodels.FileInCommitSet, error)
	GetCommitsByCommitHash(params gitmodels.CommitRequest, requestUserID *int64) (gitmodels.CommitSet, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64, requestUserID *int64) (gitmodels.CommitSet, error)
	GetBranchInfoByNames(userLogin, repoName, branchName string, currUserID *int64) (gitmodels.Branch, error)
	GetFileByPath(params gitmodels.FilesCommitRequest, requestUserID *int64) (file gitmodels.FileCommitted, err error)
	GetFileContentByBranch(userLogin, repoName, branchName, filePath string, currUserID *int64) ([]byte, error)
	GetFileContentByCommitHash(userLogin, repoName, commitHash, filePath string, currUserID *int64) ([]byte, error)
	GetRepoHead(userLogin, repoName string, requestUserID *int64) (gitmodels.Branch, error)
	Fork(repoID int64, author, repoName, newName string, currUserID int64) error
}
