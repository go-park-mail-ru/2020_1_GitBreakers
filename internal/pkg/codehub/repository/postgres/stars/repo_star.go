package stars

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type StarRepository struct {
	DB *sqlx.DB
}

func NewStarRepository(db *sqlx.DB) StarRepository {
	return StarRepository{
		DB: db,
	}
}

func (repo *StarRepository) IsExistStar(userID int64, repoID int64) (bool, error) {
	var isExist bool
	err := repo.DB.QueryRow(`
		SELECT EXISTS(
               SELECT *
               FROM git_repository_user_stars
               WHERE repository_id = $1
                 AND author_id = $2
           )`, repoID, userID).Scan(&isExist)

	if err != nil {
		return false, errors.Wrapf(err, "error occurs in StarRepository in IsExistStar function "+
			"with userId=%v, repoID=%v", userID, repoID)
	}
	return isExist, nil
}

func (repo *StarRepository) AddStar(userID int64, repoID int64) error {
	isExist, err := repo.IsExistStar(userID, repoID)
	if err != nil {
		return errors.Wrapf(err, "error occurs in StarRepository in AddStar function "+
			"with userId=%v, repoID=%v", repoID, userID)
	}

	if isExist {
		return entityerrors.AlreadyExist()
	}

	_, err = repo.DB.Exec(
		"INSERT INTO git_repository_user_stars (repository_id, author_id) VALUES ($1, $2)",
		repoID, userID)
	if err != nil {
		return errors.Wrapf(err, "error occurs in StarRepository in AddStar function "+
			"with userId=%v, repoID=%v", userID, repoID)
	}

	return nil
}

func (repo *StarRepository) DelStar(userID int64, repoID int64) error {
	var isDeleted bool

	err := repo.DB.QueryRow(
		`DELETE FROM git_repository_user_stars
				WHERE repository_id = $1 AND author_id = $2
				RETURNING TRUE AS result`,
		repoID, userID).Scan(&isDeleted)

	switch {
	case err == sql.ErrNoRows:
		return entityerrors.DoesNotExist()
	case err != nil:
		return errors.Wrapf(err, "error occurs in StarRepository in DelStar function "+
			"with userId=%v, repoID=%v", userID, repoID)
	}

	return nil
}

func (repo *StarRepository) GetStarredRepos(userID int64, limit int64, offset int64) (gitRepos []gitmodels.Repository, err error) {
	rows, err := repo.DB.Query(
		`	SELECT 	gruv.id,
					   	gruv.owner_id,
					   	gruv.name,
					   	gruv.description,
					   	gruv.is_fork,
					   	gruv.created_at,
					   	gruv.is_public,
	   					gruv.stars,
	       				gruv.user_login
				FROM git_repository_user_stars AS grus
						JOIN git_repository_user_view AS gruv ON gruv.id = grus.repository_id
				WHERE grus.author_id = $1 LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "error occurs in StarRepository in GetStarredRepos function "+
			"with userId=%v", userID)
	}

	defer func() {
		if errRows := rows.Close(); errRows != nil {
			err = errors.Wrap(err, errRows.Error())
		}
	}()

	for rows.Next() {
		gitRepo := gitmodels.Repository{}
		err = rows.Scan(
			&gitRepo.ID,
			&gitRepo.OwnerID,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.CreatedAt,
			&gitRepo.IsPublic,
			&gitRepo.Stars,
			&gitRepo.AuthorLogin,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in StarRepository in GetStarredRepos function "+
				"while scanning repositories with starUserId=%v", userID)
		}

		gitRepos = append(gitRepos, gitRepo)
	}

	return gitRepos, nil
}
func (repo *StarRepository) GetUserStaredList(repoID int64, limit int64, offset int64) (users []models.User, err error) {
	rows, err := repo.DB.Query(`
				SELECT 	upv.id,
					   	upv.login,
					   	upv.email,
					   	upv.name,
					   	upv.avatar_path,
					   	upv.created_at
				FROM git_repository_user_stars AS grus
						 JOIN user_profile_view AS upv ON grus.author_id = upv.id
				WHERE grus.repository_id = $1
				LIMIT $2 OFFSET $3`,
		repoID, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "error occurs in StarRepository in GetUserStaredList function "+
			"with repoID=%v", repoID)
	}

	defer func() {
		if errRows := rows.Close(); errRows != nil {
			err = errors.Wrap(err, errRows.Error())
		}
	}()

	for rows.Next() {
		user := models.User{}
		err = rows.Scan(
			&user.ID,
			&user.Login,
			&user.Email,
			&user.Name,
			&user.Image,
			&user.CreatedAt)
		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in StarRepository in GetUserStaredList function "+
				"while scanning users with repoID=%v", repoID)
		}
		users = append(users, user)
	}

	return users, nil
}
