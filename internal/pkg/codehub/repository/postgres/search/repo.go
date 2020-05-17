package search

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/jmoiron/sqlx"
)

type RepoSearch struct {
	DB *sqlx.DB
}

func NewSearchRepository(db *sqlx.DB) RepoSearch {
	return RepoSearch{
		DB: db,
	}
}

func (repo RepoSearch) GetFromUsers(query string, limit int64, offset int64) (models.UserSet, error) {
	var userlist []models.User
	err := repo.DB.Select(&userlist, `select id, login, email, password, name, avatar_path, created_at from users where login like $1 || '%' limit $2 offset $3`,
		query, limit, offset)

	switch {
	case err == sql.ErrNoRows:
		break
	default:
		return userlist, err
	}

	return userlist, nil
}

func (repo RepoSearch) GetFromStarredRepos(query string, limit int64, offset int64, userID int64) (models.RepoSet, error) {
	var repolist []gitmodels.Repository
	err := repo.DB.Select(&repolist, `select * from git_repositories where id in 
					  (select repository_id from git_repository_user_stars where
				  		author_id=$1) and name like $4 ||'%' order by created_at limit $2 offset $3`,
		userID, limit, offset, query)

	switch {
	case err == sql.ErrNoRows:
		break
	default:
		return repolist, err
	}

	return repolist, nil
}

func (repo RepoSearch) GetFromAllRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	var repolist []gitmodels.Repository
	err := repo.DB.Select(&repolist, `select * from git_repositories where name like $1 || '%' limit $2 offset $3`,
		query, limit, offset)

	switch {
	case err == sql.ErrNoRows:
		break
	default:
		return repolist, err
	}

	return repolist, nil
}

func (repo RepoSearch) GetFromOwnRepos(query string, limit int64, offset int64, userID int64) (models.RepoSet, error) {
	var repolist []gitmodels.Repository
	err := repo.DB.Select(&repolist, `select * from git_repositories  join 
    	users_git_repositories us_git on id=repository_id where
				  		user_id=$1 and name like $4 ||'%' order by us_git.created_at limit $2 offset $3`,
		userID, limit, offset, query)

	switch {
	case err == sql.ErrNoRows:
		break
	default:
		return repolist, err
	}

	return repolist, nil
}
