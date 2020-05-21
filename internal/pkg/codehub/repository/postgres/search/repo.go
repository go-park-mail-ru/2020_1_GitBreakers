package search

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
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
	userlist := []models.User{}
	err := repo.DB.Select(&userlist, `select * from user_profile_view where login like $1 || '%' limit $2 offset $3`,
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
	rows, err := repo.DB.Query(`
		SELECT 	repo.id,
				repo.owner_id,
				repo.name,
				repo.description,
				repo.is_fork,
				repo.is_public,
		       	repo.stars,
		       	repo.forks,
		       	repo.created_at,
		       	repo.user_login,
    			repo.parent_id,
				repo.parent_owner_id,
				repo.parent_name,
				repo.parent_user_login
		FROM git_repository_parent_user_view AS repo
			WHERE repo.id in 
					  (select repository_id from git_repository_user_stars where
				  		author_id=$1) and name like $4 ||'%' order by created_at limit $2 offset $3`,
		userID, limit, offset, query)

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	var gitRepos []gitmodels.Repository
	for rows.Next() {
		gitRepo := gitmodels.Repository{}
		err = rows.Scan(
			&gitRepo.ID,
			&gitRepo.OwnerID,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.IsPublic,
			&gitRepo.Stars,
			&gitRepo.Forks,
			&gitRepo.CreatedAt,
			&gitRepo.AuthorLogin,
			&gitRepo.ParentRepositoryInfo.ID,
			&gitRepo.ParentRepositoryInfo.OwnerID,
			&gitRepo.ParentRepositoryInfo.Name,
			&gitRepo.ParentRepositoryInfo.AuthorLogin,
		)
		if err != nil {
			return nil, err
		}

		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}

func (repo RepoSearch) GetFromAllRepos(query string, limit int64, offset int64) (models.RepoSet, error) {
	rows, err := repo.DB.Query(`
		SELECT 	repo.id,
				repo.owner_id,
				repo.name,
				repo.description,
				repo.is_fork,
				repo.is_public,
		       	repo.stars,
		       	repo.forks,
		       	repo.created_at,
		       	repo.user_login,
    			repo.parent_id,
				repo.parent_owner_id,
				repo.parent_name,
				repo.parent_user_login
		FROM git_repository_parent_user_view AS repo
		 where repo.name like $1 || '%' limit $2 offset $3`,
		query, limit, offset)

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	var gitRepos []gitmodels.Repository
	for rows.Next() {
		gitRepo := gitmodels.Repository{}
		err = rows.Scan(
			&gitRepo.ID,
			&gitRepo.OwnerID,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.IsPublic,
			&gitRepo.Stars,
			&gitRepo.Forks,
			&gitRepo.CreatedAt,
			&gitRepo.AuthorLogin,
			&gitRepo.ParentRepositoryInfo.ID,
			&gitRepo.ParentRepositoryInfo.OwnerID,
			&gitRepo.ParentRepositoryInfo.Name,
			&gitRepo.ParentRepositoryInfo.AuthorLogin,
		)
		if err != nil {
			return nil, err
		}

		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}

func (repo RepoSearch) GetFromOwnRepos(query string, limit int64, offset int64, userID int64) (models.RepoSet, error) {
	rows, err := repo.DB.Query(`SELECT 	repo.id,
				repo.owner_id,
				repo.name,
				repo.description,
				repo.is_fork,
				repo.is_public,
		       	repo.stars,
		       	repo.forks,
		       	repo.created_at,
		       	repo.user_login,
    			repo.parent_id,
				repo.parent_owner_id,
				repo.parent_name,
				repo.parent_user_login
		FROM git_repository_parent_user_view AS repo  join 
    	users_git_repositories us_git on repo.id=repository_id where
				  		user_id=$1 and name like $4 ||'%' order by us_git.created_at limit $2 offset $3`,
		userID, limit, offset, query)

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	var gitRepos []gitmodels.Repository
	for rows.Next() {
		gitRepo := gitmodels.Repository{}
		err = rows.Scan(
			&gitRepo.ID,
			&gitRepo.OwnerID,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.IsPublic,
			&gitRepo.Stars,
			&gitRepo.Forks,
			&gitRepo.CreatedAt,
			&gitRepo.AuthorLogin,
			&gitRepo.ParentRepositoryInfo.ID,
			&gitRepo.ParentRepositoryInfo.OwnerID,
			&gitRepo.ParentRepositoryInfo.Name,
			&gitRepo.ParentRepositoryInfo.AuthorLogin,
		)
		if err != nil {
			return nil, err
		}

		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}
