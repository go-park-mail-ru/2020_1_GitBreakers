package usecase

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"

type GitUseCase struct {
	repo git.Repository
}

func (GU *GitUseCase) GetBranchList(repoName string, userName string) {

}
