package news

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
	news := []models.News{}
	err := R.DB.Select(&news,
		`	SELECT	id,
	       				author_id,
	       				repository_id,
	       				message,
	       				label,
	       				created_at,
	       				user_login
				FROM news_users_view WHERE repository_id=$1 LIMIT $2 OFFSET $3`,
		repoID, limit, offset)

	return news, err
}
