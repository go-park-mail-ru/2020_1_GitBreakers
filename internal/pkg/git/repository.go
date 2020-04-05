package git

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type Repository interface {
	GetById(id int) (git.Repository, error) // TODO
	GetByName(userLogin, repoName string) (git.Repository, error)
	Create(repos git.Repository) error
	DeleteById(id int) error // TODO
	DeleteByName(userId int, repoName string) error
	CheckReadAccess(currentUserId *int, userLogin, repoName string) (bool, error)
	GetBranchesByName(userLogin, repoName string) ([]git.Branch, error)
	GetAnyReposByUserLogin(userLogin string, offset, limit int) ([]git.Repository, error)
	GetReposByUserLogin(requesterId *int, userLogin string, offset, limit int) ([]git.Repository, error)

	GetCommitsInBranch(repoName, branchName string, offset, limit int) ([]git.Commit, error)
	FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error)
}
