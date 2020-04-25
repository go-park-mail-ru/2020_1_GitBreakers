package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoNewsI interface {
	GetNews(repoID int64) (models.News, error)
}
