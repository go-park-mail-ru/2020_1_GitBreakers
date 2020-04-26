package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoNewsI interface {
	GetNews(repoID int64, limit int64, offset int64) (models.NewsSet, error)
}
