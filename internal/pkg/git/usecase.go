package git

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type GitUseCaseI interface {
	Create(userID int64, repos *gitmodels.Repository) error
	GetRepo(userName string, repoName string, requestUserID *int64) (gitmodels.Repository, error)
	GetRepoList(userName string, requestUserID *int64) (gitmodels.RepositorySet, error)
	GetBranchList(requestUserID *int64, userName string, repoName string) (gitmodels.BranchSet, error)
	FilesInCommitByPath(requset gitmodels.FilesCommitRequest, requesrUserID *int64) (gitmodels.FileInCommitSet, error)
	GetCommitsByCommitHash(params gitmodels.CommitRequest, requestUserID *int64) (gitmodels.CommitSet, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64, requestUserID *int64) (gitmodels.CommitSet, error)
	GetFileByPath(params gitmodels.FilesCommitRequest, requestUserID *int64) (file gitmodels.FileCommitted, err error)
	GetRepoHead(userLogin, repoName string, requestUserID *int64) (gitmodels.Branch, error)
	Fork(repoID int64, author, repoName, newName string, currUserID int64) error
}
