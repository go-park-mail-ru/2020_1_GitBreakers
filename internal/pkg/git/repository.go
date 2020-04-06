package git

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type Repository interface {
	GetById(id int) (git.Repository, error)
	GetByName(userLogin, repoName string) (git.Repository, error)
	Create(repos git.Repository) (id int64, err error)
	//DeleteById(id int) error // TODO
	//DeleteByName(userId int, repoName string) error // TODO
	CheckReadAccess(currentUserId *int, userLogin, repoName string) (bool, error)
	GetBranchesByName(userLogin, repoName string) ([]git.Branch, error)
	GetAnyReposByUserLogin(userLogin string, offset, limit int) ([]git.Repository, error)
	GetReposByUserLogin(requesterId *int, userLogin string, offset, limit int) ([]git.Repository, error)

	FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error)
	GetFileByPath(userLogin, repoName, commitHash, path string) (git.FileCommitted, error)
	GetCommitsByCommitHash(userLogin, repoName, commitHash string, offset, limit int) ([]git.Commit, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int) ([]git.Commit, error)
}
