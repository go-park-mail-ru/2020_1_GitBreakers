package git

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type UseCase interface {
	Create(userID int64, repos *gitmodels.Repository) error
	GetRepo(userName string, repoName string, requestUserID *int64) (gitmodels.Repository, error)
	GetRepoList(userName string, requestUserID *int64) ([]gitmodels.Repository, error)
	GetBranchList(requestUserID *int64, userName string, repoName string) ([]gitmodels.Branch, error)
	FilesInCommitByPath(requset gitmodels.FilesCommitRequest, requesrUserID *int64) ([]gitmodels.FileInCommit, error)
	GetCommitsByCommitHash(params gitmodels.CommitRequest, requestUserID *int64) ([]gitmodels.Commit, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64, requestUserID *int64) ([]gitmodels.Commit, error)
	GetFileByPath(params gitmodels.FilesCommitRequest, requestUserID *int64) (file gitmodels.FileCommitted, err error)
}
