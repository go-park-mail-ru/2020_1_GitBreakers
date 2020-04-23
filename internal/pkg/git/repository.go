package git

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type Repository interface {
	GetByID(id int64) (git.Repository, error)
	GetByName(userLogin, repoName string) (git.Repository, error)
	Create(repos git.Repository) (id int64, err error)
	//DeleteByID(id int) error // TODO
	//DeleteByName(userId int, repoName string) error // TODO
	CheckReadAccess(currentUserId *int64, userLogin, repoName string) (bool, error)
	GetBranchesByName(userLogin, repoName string) ([]git.Branch, error)
	GetAnyReposByUserLogin(userLogin string, offset, limit int64) ([]git.Repository, error)
	GetReposByUserLogin(requesterId *int64, userLogin string, offset, limit int64) ([]git.Repository, error)

	FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error)
	GetFileByPath(userLogin, repoName, commitHash, path string) (git.FileCommitted, error)
	GetCommitsByCommitHash(userLogin, repoName, commitHash string, offset, limit int64) ([]git.Commit, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64) ([]git.Commit, error)
}
