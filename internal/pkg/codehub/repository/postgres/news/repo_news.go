package news

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	var news models.NewsSet

	rows, err := R.DB.Query(
		`	SELECT	id,
	       				author_id,
	       				repository_id,
	       				message,
	       				label,
	       				created_at,
	       				user_login,
	       				user_avatar_path
				FROM news_users_view WHERE repository_id=$1 LIMIT $2 OFFSET $3`,
		repoID, limit, offset)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for rows.Next() {
		var newsModel models.News

		err := rows.Scan(
			&newsModel.ID,
			&newsModel.AuthorID,
			&newsModel.RepoID,
			&newsModel.Mess,
			&newsModel.Label,
			&newsModel.Date,
			&newsModel.AuthorLogin,
			&newsModel.AuthorImage,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		news = append(news, newsModel)
	}

	return news, err
}
