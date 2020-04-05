package git

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type Repository interface {
	GetById(id int) (git.Repository, error) // TODO
	GetByName(userLogin, repoName string) (git.Repository, error) // DONE
	Create(repos git.Repository) (id int64, err error) // DONE
	DeleteById(id int) error // TODO
	DeleteByName(userId int, repoName string) error // TODO
	CheckReadAccess(currentUserId *int, userLogin, repoName string) (bool, error)  // DONE
	GetBranchesByName(userLogin, repoName string) ([]git.Branch, error)  // DONE
	GetAnyReposByUserLogin(userLogin string, offset, limit int) ([]git.Repository, error)  // DONE
	GetReposByUserLogin(requesterId *int, userLogin string, offset, limit int) ([]git.Repository, error)  // DONE

	GetCommitsInBranch(repoName, branchName string, offset, limit int) ([]git.Commit, error) // TODO
	FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error) // TODO
}
