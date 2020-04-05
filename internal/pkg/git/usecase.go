package git

import gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"

type UseCase interface {
	Create(userid int, repos *gitmodels.Repository) error
	//Update()
	GetRepo(userName string, repoName string) (gitmodels.Repository, error)
	GetRepoList(userName string) ([]gitmodels.Repository, error)
	//GetBranchList(userName string, repoName string) gitmodels.Branch
}
