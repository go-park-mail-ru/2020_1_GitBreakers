package gitserver

import (
	gitModels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type UseCase interface {
	IsGitRepositoryExists(userLogin string, repoName string) (bool, error)
	GetGitRepository(userLogin string, repoName string) (gitModels.Repository, error)
	CheckUserPassword(userLogin string, password string) (bool, error)
	CheckGitRepositoryReadAccess(currentUserId *int64, userLogin, repoName string) (bool, error)
	GetGitRepositoryPermission(currentUserId *int64, userLogin, repoName string) (perm.Permission, error)
}
