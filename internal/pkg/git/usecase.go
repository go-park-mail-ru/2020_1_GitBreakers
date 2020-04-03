package git

import gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"

type UseCase interface {
	Create(userid int, repository *gitmodels.Repository)
	Update()
	GetRepo(userName string, repoName string) gitmodels.Repository
	GetRepoList(userName string) *[]gitmodels.Repository
}
