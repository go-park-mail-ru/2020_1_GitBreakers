package git

import (
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type UseCase interface {
	Create(userid int, repos *gitmodels.Repository) error
	//Update()
	GetRepo(userName string, repoName string) (gitmodels.Repository, error)
	GetRepoList(userName string) ([]gitmodels.Repository, error)
	GetBranchList(requestUserID *int, userName string, repoName string) ([]gitmodels.Branch, error)
	FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]gitmodels.FileInCommit, error)
	GetCommitsByCommitHash(params gitmodels.CommitRequest) ([]gitmodels.Commit, error)

}
