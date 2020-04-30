package interfaces

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type NewsClientI interface {
	GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error)
}
