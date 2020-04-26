package postgres

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type RepoNews struct {
	DB *sqlx.DB
}

func NewRepoNews(db *sqlx.DB) codehub.RepoNewsI {
	return &RepoNews{DB: db}
}
func (R *RepoNews) GetNews(repoID int64, limit int64, offset int64) (models.NewsSet, error) {
	panic("implement me")
}
