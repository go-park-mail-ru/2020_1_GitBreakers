package git

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type GitRepoI interface {
	GetByID(id int64) (git.Repository, error)
	GetByName(userLogin, repoName string) (git.Repository, error)
	Create(repos git.Repository) (id int64, err error)
	DeleteByOwnerID(ownerID int64, repoName string) error

	CheckReadAccess(currentUserId *int64, userLogin, repoName string) (bool, error)
	CheckReadAccessById(currentUserId *int64, repoId int64) (bool, error)
	// GetPermission returns permission: for public repo - write and higher, for private - read and higher
	// In other case returns NoAccess
	GetPermission(currentUserId *int64, userLogin, repoName string) (perm.Permission, error)
	GetPermissionByID(currentUserId *int64, repoID int64) (perm.Permission, error)
	GetRepoPathByID(repoID int64) (string, error)

	IsRepoExistsByID(repoID int64) (bool, error)
	IsRepoExistsByOwnerId(ownerId int64, repoName string) (bool, error)
	IsRepoExistsByOwnerLogin(ownerLogin string, repoName string) (bool, error)
	GetBranchHashIfExistInRepoByID(repoID int64, branchName string) (string, error)

	GetBranchesByName(userLogin, repoName string) ([]git.Branch, error)
	GetAnyReposByUserLogin(userLogin string, offset, limit int64) ([]git.Repository, error)
	GetReposByUserLogin(requesterId *int64, userLogin string, offset, limit int64) ([]git.Repository, error)

	FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error)
	GetFileByPath(userLogin, repoName, commitHash, path string) (git.FileCommitted, error)
	GetCommitsByCommitHash(userLogin, repoName, commitHash string, offset, limit int64) ([]git.Commit, error)
	GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64) ([]git.Commit, error)
	GetRepoHead(userLogin, repoName string) (git.Branch, error)

	Fork(forkRepoName string, userID, repoBaseID int64) error
}
