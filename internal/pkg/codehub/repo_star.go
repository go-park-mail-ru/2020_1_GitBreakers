package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type RepoStarI interface {
	IsExistStar(userID int64, repoID int64) (bool, error)
	AddStar(userID int64, repoID int64) error
	DelStar(userID int64, repoID int64) error
	GetStarredRepos(userID int64, limit int64, offset int64) ([]gitmodels.Repository, error)
	GetUserStaredList(repoID int64, limit int64, offset int64) ([]models.User, error)
}
