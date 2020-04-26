package codehub

import gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"

type RepoStarI interface {
	AddStar(userID int64, repoID int64) error
	DelStar(userID int64, repoID int64) error
	GetStarredRepo(userID int64, limit int64, offset int64) ([]gitmodels.Repository, error)
}


