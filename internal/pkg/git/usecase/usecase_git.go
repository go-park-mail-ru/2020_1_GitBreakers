package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type GitUseCase struct {
	Repo git.Repository
}

//func (GU *GitUseCase) GetBranchList(repoName string, userName string) {
//
//}
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
