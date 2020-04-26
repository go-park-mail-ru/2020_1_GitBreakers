package postgres

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
               FROM git_repository_user_star
               WHERE repository_id = $1
                 AND user_id = $2
           )`, repoID, userID).Scan(&isExist)

	if err != nil {
		return false, errors.Wrapf(err, "error occurs in StarRepository in IsExistStar function "+
			"with userId=%v, repoID=%v", userID, repoID)
	}
	return true, nil
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
		"INSERT INTO git_repository_user_star (repository_id, user_id) VALUES ($1, $2)",
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
		"DELETE FROM git_repository_user_star WHERE repository_id = $1 AND user_id = $2 RETURNING TRUE",
		repoID, userID).Scan(&isDeleted)

	switch true {
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
		`	SELECT gr.id,
					   gr.owner_id,
					   gr.name,
					   gr.description,
					   gr.is_fork,
					   gr.created_at,
					   gr.is_public
				FROM git_repository_user_star AS grus
						 JOIN git_repositories AS gr ON grus.repository_id = gr.id
				WHERE grus.user_id = $1 LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, errors.Wrapf(err, "error occurs in StarRepository in GetStarredRepos function "+
			"with userId=%v", userID)
	}

	defer func() {
		errRows := rows.Close()
		if errRows != nil {
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
			&gitRepo.IsPublic)
		if err != nil {
			return nil, errors.Wrapf(err, "error occurs in StarRepository in GetStarredRepos function "+
				"while scanning repositories with starUserId=%v", userID)
		}
		gitRepos = append(gitRepos, gitRepo)
	}

	return gitRepos, nil
}
func (repo *StarRepository) GetUserStaredList(repoID int64, limit int64, offset int64) ([]models.User, error) {
	panic("implement me")
}
