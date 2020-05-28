package search

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
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
	userList := make(models.UserSet, 0)

	rows, err := repo.DB.Query(`
			SELECT id,
				   login,
				   email,
				   name,
				   avatar_path,
				   created_at
			FROM user_profile_view
			WHERE login LIKE $1 || '%' ORDER BY name LIMIT $2
			OFFSET $3`,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			if err == nil {
				err = closeErr
			} else {
				err = errors.Wrap(err, closeErr.Error())
			}
		}
	}()

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.Name,
			&user.Image,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		userList = append(userList, user)
	}

	return userList, nil
}

func (repo RepoSearch) GetFromStarredRepos(query string, limit int64, offset int64, userID *int64) (models.RepoSet, error) {
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
					JOIN git_repository_user_stars AS grus ON repo.id = grus.repository_id 
				AND grus.author_id = $1
					JOIN users_git_repositories AS ugr ON repo.id = ugr.repository_id
				AND (ugr.user_id = $1 AND ugr.role NOT IN ($5) OR repo.is_public = TRUE)
			WHERE name LIKE $4 || '%' ORDER BY repo.created_at LIMIT $2 
			OFFSET $3`,
		userID,
		limit,
		offset,
		query,
		perm.NoAccess(),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	gitRepos, err := scanRepositories(rows)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return gitRepos, nil
}

func (repo RepoSearch) GetFromAllRepos(query string, limit int64, offset int64, userID *int64) (models.RepoSet, error) {
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
					JOIN users_git_repositories AS ugr ON repo.id = ugr.repository_id
				AND (ugr.user_id = $1 AND ugr.role NOT IN ($5) OR repo.is_public = TRUE)
			WHERE repo.name LIKE $2 || '%' ORDER BY repo.created_at LIMIT $3
			OFFSET $4`,
		userID,
		query,
		limit,
		offset,
		perm.NoAccess(),
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	gitRepos, err := scanRepositories(rows)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return gitRepos, nil
}

func (repo RepoSearch) GetFromOwnRepos(query string, limit int64, offset int64, userID *int64) (models.RepoSet, error) {
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
		    	JOIN users_git_repositories AS ugr ON repo.id = ugr.repository_id 
			AND ugr.user_id = $1 AND ugr.role IN ($5) 
		WHERE name LIKE $4 ||'%' ORDER BY repo.created_at LIMIT $2
		OFFSET $3`,
		userID, limit, offset, query, perm.OwnerAccess())

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	gitRepos, err := scanRepositories(rows)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return gitRepos, nil
}

func scanRepositories(rows *sql.Rows) (models.RepoSet, error) {
	gitRepos := make(models.RepoSet, 0)

	for rows.Next() {
		gitRepo := gitmodels.Repository{}
		err := rows.Scan(
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
			return nil, errors.WithStack(err)
		}

		gitRepos = append(gitRepos, gitRepo)
	}

	return gitRepos, nil
}
