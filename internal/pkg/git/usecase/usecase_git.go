package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/pkg/errors"
)

type GitUseCase struct {
	Repo git.Repository
}

func (GU *GitUseCase) Create(userid int, repos *gitmodels.Repository) error {
	repos.OwnerId = userid
	_, err := GU.Repo.Create(*repos)
	switch {
	case err == entityerrors.AlreadyExist():
		return entityerrors.AlreadyExist()
	case err != nil:
		return errors.Wrap(err, "repo not created")
	}
	return nil
}

func (GU *GitUseCase) GetRepo(userName string, repoName string) (gitmodels.Repository, error) {
	return GU.Repo.GetByName(userName, repoName)
}
func (GU *GitUseCase) GetRepoList(userName string) ([]gitmodels.Repository, error) {
	return GU.Repo.GetAnyReposByUserLogin(userName, 0, 100)
}

func (GU *GitUseCase) GetBranchList(requestUserID *int, userName string, repoName string) ([]gitmodels.Branch, error) {
	ReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, userName, repoName)
	if ReadyToRead {
		return GU.Repo.GetBranchesByName(userName, repoName)
	} else {
		return nil, err
	}
}
func (GU *GitUseCase) FilesInCommitByPath(request gitmodels.FilesCommitRequest) ([]gitmodels.FileInCommit, error) {
	return GU.Repo.FilesInCommitByPath(request.UserName, request.Reponame, request.HashCommits, request.Path)
}
func (GU *GitUseCase) GetCommitsByCommitHash(params gitmodels.CommitRequest) ([]gitmodels.Commit, error) {
	if params.Limit == 0 {
		params.Limit = 100
	}
	return GU.Repo.GetCommitsByCommitHash(params.UserLogin,
		params.RepoName, params.CommitHash, params.Offset, params.Limit)
}
func (GU *GitUseCase) GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int) ([]gitmodels.Commit, error) {
	return GU.Repo.GetCommitsByBranchName(userLogin, repoName, branchName, offset, limit)
}
func (GU *GitUseCase) GetFileByPath(params gitmodels.FilesCommitRequest) (file gitmodels.FileCommitted, err error) {
	return GU.Repo.GetFileByPath(params.UserName, params.Reponame, params.HashCommits, params.Path)
}
