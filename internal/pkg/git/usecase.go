package git

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type UseCase interface {
	Create(userid int, repos *gitmodels.Repository) error
	GetRepo(userName string, repoName string, requestUserID *int) (gitmodels.Repository, error)
	GetRepoList(userName string, requestUserID *int) ([]gitmodels.Repository, error)
	GetBranchList(requestUserID *int, userName string, repoName string) ([]gitmodels.Branch, error)
	FilesInCommitByPath(requset gitmodels.FilesCommitRequest, requesrUserID *int) ([]gitmodels.FileInCommit, error)
	GetCommitsByCommitHash(params gitmodels.CommitRequest, requestUserID *int) ([]gitmodels.Commit, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int, requestUserID *int) ([]gitmodels.Commit, error)
	GetFileByPath(params gitmodels.FilesCommitRequest, requestUserID *int) (file gitmodels.FileCommitted, err error)
}
