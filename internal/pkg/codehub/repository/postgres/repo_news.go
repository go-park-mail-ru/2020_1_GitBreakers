package postgres

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type RepoNews struct {
	DB *sqlx.DB
}

func NewRepoNews(db *sqlx.DB) RepoNews {
	return RepoNews{
		DB: db,
	}
}

func (R *RepoNews) GetNews(repoID int64, limit int64, offset int64) (models.NewsSet, error) {
	news := models.NewsSet{}
	err := R.DB.Select(&news,
		`select * from news where repo_id=$1 limit $2 offset $3`,
		repoID, limit, offset)

	return news, err
}
