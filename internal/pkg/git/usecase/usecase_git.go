package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type GitUseCase struct {
	repo git.Repository
}

func (GU *GitUseCase) GetBranchList(repoName string, userName string) {

}
func (GU *GitUseCase) Create(userid int, repos *gitmodels.Repository) {
	//todo лучше бы по указателю принимал
	repos.OwnerId = userid
	GU.repo.Create(*repos)

}
func (GU *GitUseCase) GetRepoByName(repoName string, userName string) {
	GU.repo.GetByName(userName, repoName)
}
