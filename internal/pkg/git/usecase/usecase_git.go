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

func (GU *GitUseCase) GetRepo(userName string, repoName string, requestUserID *int) (gitmodels.Repository, error) {
	isReadAccepted, err := GU.Repo.CheckReadAccess(requestUserID, userName, repoName)
	if err != nil {
		return gitmodels.Repository{}, errors.Wrap(err, "error in access check")
	}

	if !isReadAccepted {
		return gitmodels.Repository{}, entityerrors.AccessDenied()
	}

	return GU.Repo.GetByName(userName, repoName)
}
func (GU *GitUseCase) GetRepoList(userName string, requestUserID *int) ([]gitmodels.Repository, error) {
	rawRepoList, err := GU.Repo.GetAnyReposByUserLogin(userName, 0, 100)
	if err != nil {
		return nil, errors.Wrap(err, "didn't get repolist")
	}

	resultRepoList := make([]gitmodels.Repository, 0)
	for _, v := range rawRepoList {
		isReadAccepted, err := GU.Repo.CheckReadAccess(requestUserID, userName, v.Name)
		if err == nil && isReadAccepted {
			resultRepoList = append(resultRepoList, v)
		}
	}

	return resultRepoList, nil
}

func (GU *GitUseCase) GetBranchList(requestUserID *int, userName string, repoName string) ([]gitmodels.Branch, error) {
	ReadyToRead, _ := GU.Repo.CheckReadAccess(requestUserID, userName, repoName)

	if ReadyToRead {
		return GU.Repo.GetBranchesByName(userName, repoName)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) FilesInCommitByPath(request gitmodels.FilesCommitRequest, requestUserID *int) ([]gitmodels.FileInCommit, error) {
	ReadyToRead, _ := GU.Repo.CheckReadAccess(requestUserID, request.UserName, request.Reponame)

	if ReadyToRead {
		return GU.Repo.FilesInCommitByPath(request.UserName, request.Reponame, request.HashCommits, request.Path)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) GetCommitsByCommitHash(params gitmodels.CommitRequest, requestUserID *int) ([]gitmodels.Commit, error) {
	ReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, params.UserLogin, params.RepoName)

	if ReadyToRead && err == nil {
		if params.Limit == 0 {
			params.Limit = 100
		}
		return GU.Repo.GetCommitsByCommitHash(params.UserLogin,
			params.RepoName, params.CommitHash, params.Offset, params.Limit)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int, requestUserID *int) ([]gitmodels.Commit, error) {
	ReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, userLogin, repoName)

	if ReadyToRead && err == nil {
		return GU.Repo.GetCommitsByBranchName(userLogin, repoName, branchName, offset, limit)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) GetFileByPath(params gitmodels.FilesCommitRequest, requestUserID *int) (file gitmodels.FileCommitted, err error) {
	ReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, params.UserName, params.Reponame)

	if ReadyToRead && err == nil {
		return GU.Repo.GetFileByPath(params.UserName, params.Reponame, params.HashCommits, params.Path)
	}

	return gitmodels.FileCommitted{}, entityerrors.AccessDenied()
}
