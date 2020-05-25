package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitModels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/pkg/errors"
)

type UseCase struct {
	gitRepoRepository git.GitRepoI
	mergeRepo         codehub.RepoMergeI
	userClient        clients.UserClient
}

func NewUseCase(gitRepo git.GitRepoI, mergeRepo codehub.RepoMergeI, client clients.UserClient) UseCase {
	return UseCase{
		gitRepoRepository: gitRepo,
		mergeRepo:         mergeRepo,
		userClient:        client,
	}
}

func (u UseCase) IsGitRepositoryExists(userLogin string, repoName string) (bool, error) {
	_, err := u.gitRepoRepository.GetByName(userLogin, repoName)
	switch true {
	case err == entityerrors.ErrDoesNotExist:
		return false, nil
	case err != nil:
		return false, errors.Wrapf(err, "error while checking if git repository exists: "+
			"userLogin: %v, repoName=%v", userLogin, repoName)
	}
	return true, nil
}

func (u UseCase) GetGitRepository(userLogin string, repoName string) (gitModels.Repository, error) {
	repo, err := u.gitRepoRepository.GetByName(userLogin, repoName)
	if err != nil {
		return repo, errors.Wrapf(err, "error while getting git repository: "+
			"userLogin: %v, repoName=%v", userLogin, repoName)
	}
	return repo, nil
}

func (u UseCase) GetUser(userLogin string) (models.User, error) {
	user, err := u.userClient.GetByLogin(userLogin)
	if err != nil {
		return user, errors.Wrapf(err, "error while getting user by userLogin=%v", userLogin)
	}
	return user, nil
}

func (u UseCase) CheckUserPassword(userLogin string, password string) (bool, error) {
	status, err := u.userClient.CheckPass(userLogin, password)
	if err != nil {
		return false, err
	}
	return status, nil
}

func (u UseCase) CheckGitRepositoryReadAccess(currentUserId *int64, userLogin, repoName string) (bool, error) {
	readAccess, err := u.gitRepoRepository.CheckReadAccess(currentUserId, userLogin, repoName)
	if err != nil {
		return false, err
	}
	return readAccess, nil
}

func (u UseCase) GetGitRepositoryPermission(currentUserId *int64, userLogin, repoName string) (perm.Permission, error) {
	permission, err := u.gitRepoRepository.GetPermission(currentUserId, userLogin, repoName)
	if err != nil {
		return permission, err
	}
	return permission, nil
}

func (u UseCase) UpdateMergeRequestsStatuses(ownerLogin, repoName string) error {
	repo, err := u.GetGitRepository(ownerLogin, repoName)
	if err != nil {
		return err
	}

	return u.mergeRepo.UpdateOpenedMergeRequestsStatusByRepoId(codehub.StatusNeedsUpdate, repo.ID)
}
