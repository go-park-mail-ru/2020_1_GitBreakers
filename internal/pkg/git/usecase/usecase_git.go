package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type GitUseCase struct {
	Repo git.Repository
}

func (GU *GitUseCase) Create(userid int, repos *gitmodels.Repository) error {
	//todo лучше бы по указателю принимал
	repos.OwnerId = userid
	if _, err := GU.Repo.Create(*repos); err != nil {
		return err
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
	//todo настроить путь
	//if request.Path == "" {
	//	request.Path = "./"
	//}
	return GU.Repo.FilesInCommitByPath(request.UserName, request.Reponame, request.HashCommits, request.Path)
}
func (GU *GitUseCase) GetCommitsByCommitHash(params gitmodels.CommitRequest) ([]gitmodels.Commit, error) {
	if params.Limit == 0 {
		params.Limit = 100
	}
	return GU.Repo.GetCommitsByCommitHash(params.UserLogin,
		params.RepoName, params.CommitHash, params.Offset, params.Limit)
}
