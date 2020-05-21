package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoSearchI interface {
	GetFromUsers(query string, limit int64, offset int64) (models.UserSet, error)
	GetFromStarredRepos(query string, limit int64, offset int64, userID int64) (models.RepoSet, error)
	GetFromAllRepos(query string, limit int64, offset int64) (models.RepoSet, error)
	GetFromOwnRepos(query string, limit int64, offset int64, userID int64) (models.RepoSet, error)
}
