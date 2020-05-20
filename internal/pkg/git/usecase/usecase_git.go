package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/pkg/errors"
)

type GitUseCase struct {
	Repo git.GitRepoI
}

func (GU *GitUseCase) Create(userID int64, repos *gitmodels.Repository) error {
	repos.OwnerID = userID
	_, err := GU.Repo.Create(*repos)
	switch {
	case err == entityerrors.AlreadyExist():
		return entityerrors.AlreadyExist()
	case err != nil:
		return errors.Wrap(err, "repo not created")
	}
	return nil
}

func (GU *GitUseCase) GetRepo(userName string, repoName string, requestUserID *int64) (gitmodels.Repository, error) {
	isReadAccepted, err := GU.Repo.CheckReadAccess(requestUserID, userName, repoName)
	if err != nil {
		return gitmodels.Repository{}, errors.Wrap(err, "error in access check")
	}

	if !isReadAccepted {
		return gitmodels.Repository{}, entityerrors.AccessDenied()
	}

	return GU.Repo.GetByName(userName, repoName)
}

func (GU *GitUseCase) DeleteByOwnerID(ownerID int64, repoName string) error {
	// may be nil, DoesNotExist or another error
	return GU.Repo.DeleteByOwnerID(ownerID, repoName)
}

func (GU *GitUseCase) GetRepoList(userName string, requestUserID *int64) (gitmodels.RepositorySet, error) {
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

func (GU *GitUseCase) GetBranchList(requestUserID *int64, userName string, repoName string) (gitmodels.BranchSet, error) {
	ReadyToRead, _ := GU.Repo.CheckReadAccess(requestUserID, userName, repoName)

	if ReadyToRead {
		return GU.Repo.GetBranchesByName(userName, repoName)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) FilesInCommitByPath(requset gitmodels.FilesCommitRequest, requesrUserID *int64) (gitmodels.FileInCommitSet, error) {
	ReadyToRead, _ := GU.Repo.CheckReadAccess(requesrUserID, requset.UserName, requset.Reponame)

	if ReadyToRead {
		return GU.Repo.FilesInCommitByPath(requset.UserName, requset.Reponame, requset.HashCommits, requset.Path)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) GetCommitsByCommitHash(params gitmodels.CommitRequest, requestUserID *int64) (gitmodels.CommitSet, error) {
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
func (GU *GitUseCase) GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64, requestUserID *int64) (gitmodels.CommitSet, error) {
	ReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, userLogin, repoName)

	if ReadyToRead && err == nil {
		return GU.Repo.GetCommitsByBranchName(userLogin, repoName, branchName, offset, limit)
	}

	return nil, entityerrors.AccessDenied()
}
func (GU *GitUseCase) GetFileByPath(params gitmodels.FilesCommitRequest, requestUserID *int64) (file gitmodels.FileCommitted, err error) {
	ReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, params.UserName, params.Reponame)

	if ReadyToRead && err == nil {
		return GU.Repo.GetFileByPath(params.UserName, params.Reponame, params.HashCommits, params.Path)
	}

	return gitmodels.FileCommitted{}, entityerrors.AccessDenied()
}

func (GU *GitUseCase) GetRepoHead(userLogin, repoName string, requestUserID *int64) (gitmodels.Branch, error) {
	isReadyToRead, err := GU.Repo.CheckReadAccess(requestUserID, userLogin, repoName)
	if err != nil {
		return gitmodels.Branch{}, err
	}

	if !isReadyToRead {
		return gitmodels.Branch{}, entityerrors.AccessDenied()
	}

	return GU.Repo.GetRepoHead(userLogin, repoName)
}
func (GU *GitUseCase) Fork(repoID int64, author, repoName, newName string, currUserID int64) error {
	repoFromDB := gitmodels.Repository{}
	if repoID < 0 {
		var err error
		repoFromDB, err = GU.Repo.GetByName(author, repoName)
		if err != nil {
			return err
		}
		//переопределелили id
		repoID = repoFromDB.ID
	}
	isCorrectPerm, err := GU.Repo.CheckReadAccess(&currUserID, repoFromDB.AuthorLogin, repoFromDB.Name)

	if isCorrectPerm && err == nil {
		err = GU.Repo.Fork(newName, currUserID, repoID)
	}
	if err == entityerrors.DoesNotExist() {
		return err
	}
	if isCorrectPerm == false {
		return entityerrors.AccessDenied()
	}

	return err
}
