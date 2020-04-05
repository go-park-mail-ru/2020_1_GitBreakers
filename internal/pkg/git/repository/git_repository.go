package repository

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"

	gogit "github.com/go-git/go-git/v5"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	gitPostfix = ".git"
)

type Repository struct {
	db       *sqlx.DB
	reposDir string
}

func NewRepository(db *sqlx.DB, reposDirs string) Repository {
	return Repository{db: db, reposDir: reposDirs}
}

func createRepoPath(tx *sql.Tx, ownerId int, repoName string) (string, error) {
	var userLogin string
	err := tx.QueryRow("SELECT login FROM users WHERE id = $1",
		ownerId).Scan(&userLogin)
	if err != nil {
		return "", err
	}
	repoPath := userLogin + gitPostfix + "/" + repoName
	return repoPath, nil
}

func isRepoExistsInDb(tx *sql.Tx, ownerId int, repoName string) (bool, error) {
	var isRepoExists bool
	err := tx.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM git_repositories WHERE owner_id = $1 and name = $2)",
		ownerId, repoName).Scan(&isRepoExists)
	if err != nil {
		return false, errors.Wrap(err, "error while checking if git repository exists")
	}
	return isRepoExists, nil
}

func (repo Repository) Create(newRepo git.Repository) (err error) {
	// Begin transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return errors.Wrap(err, "cannot begin transaction in create repository")
	}
	// Transaction cleanups
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Wrap(err, rollbackErr.Error())
			}
		} else if commitErr := tx.Commit(); commitErr != nil {
			err = errors.Wrap(commitErr, commitErr.Error())
		}
	}()

	isRepoExist, err := isRepoExistsInDb(tx, newRepo.OwnerId, newRepo.Name)
	if err != nil {
		return errors.Wrap(err, "error in create repository while checking if repository is not exits")
	}
	if isRepoExist {
		return entityerrors.AlreadyExist()
	}

	// Create new db entity of git_repository
	repoCreationResult, err := tx.Exec(
		`INSERT INTO git_repositories (owner_id, name, description, is_public, is_fork) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		newRepo.OwnerId, newRepo.Name, newRepo.Description, newRepo.IsPublic, newRepo.IsFork)
	if err != nil {
		return errors.Wrapf(err, "cannot create new git repository entity in database, newRepo=%+v", newRepo)
	}

	_, err = tx.Exec("INSERT INTO users_git_repositories (repository_id, user_id) VALUES ($1, $2)",
		repoCreationResult.LastInsertId(), newRepo.OwnerId)
	if err != nil {
		return errors.Wrapf(err, "cannot create new git repository entity in database, newRepo=%+v", newRepo)
	}

	// Calculate path where git creates new repository on filesystem
	repoPath, err := createRepoPath(tx, newRepo.OwnerId, newRepo.Name)
	if err != nil {
		return errors.Wrapf(err, "cannot create new git repository entity in database, newRepo=%+v", newRepo)
	}

	// Create new bare repository aka 'git init --bare' on repoPath
	_, err = gogit.PlainInit(repo.reposDir+"/"+repoPath, true)
	if err == gogit.ErrRepositoryAlreadyExists {
		return entityerrors.AlreadyExist()
	}

	return nil
}

func (repo Repository) GetReposByUserLogin(requesterId *int, userLogin string, offset, limit int) ([]git.Repository, error) {
	rows, err := repo.db.Query(
		` SELECT 	repo.id,
			   			repo.owner_id,
			   			repo.name,
			   			repo.description,
			   			repo.is_fork,
			   			repo.created_at,
			   			repo.is_public
				FROM git_repositories AS repo
					JOIN users AS owner ON repo.owner_id = owner.id
					JOIN users_git_repositories AS ugr ON repo.id = ugr.repository_id
				WHERE owner.login = $1
				AND (ugr.user_id = $2 OR repo.is_public = TRUE) OFFSET $3
				LIMIT $4`,
		userLogin, requesterId, offset, limit)

	if err != nil {
		return nil, errors.Wrapf(err, "error while performing query in repository "+
			"for git repositories in GetReposByUserLogin with userLogin=%s, requesterId=%+v offset=%d, limit=$d",
			userLogin, requesterId, offset, limit)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	var gitRepos []git.Repository
	for rows.Next() {
		gitRepo := git.Repository{}
		err = rows.Scan(
			&gitRepo.Id,
			&gitRepo.OwnerId,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.CreatedAt,
			&gitRepo.IsPublic)
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories "+
				"while scanning in GetReposByUserLogin userLogin=%s, offset=%d, limit=$d", userLogin, offset, limit)
		}
		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}

func (repo Repository) GetAnyReposByUserLogin(userLogin string, offset, limit int) ([]git.Repository, error) {
	rows, err := repo.db.Query(`
		SELECT 	repo.id,
				repo.owner_id,
				repo.name,
				repo.description,
				repo.is_fork,
				repo.created_at,
				repo.is_public
		FROM git_repositories AS repo
			JOIN users AS owner ON repo.owner_id = owner.id
		WHERE owner.login = $1 OFFSET $2
		LIMIT $3`,
		userLogin, offset, limit)

	if err != nil {
		return nil, errors.Wrapf(err, "error while performing query in repository "+
			"for git repositories in GetAnyReposByUserLogin with userLogin=%s, offset=%d, limit=$d",
			userLogin, offset, limit)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	var gitRepos []git.Repository
	for rows.Next() {
		gitRepo := git.Repository{}
		err = rows.Scan(
			&gitRepo.Id,
			&gitRepo.OwnerId,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.CreatedAt,
			&gitRepo.IsPublic)
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories "+
				"while scanning in GetAnyReposByUserLogin userLogin=%s, offset=%d, limit=$d", userLogin, offset, limit)
		}
		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}

func (repo Repository) GetByName(userLogin, repoName string) (git.Repository, error) {
	var gitRepo git.Repository

	err := repo.db.QueryRow(`
		SELECT r.id,
			   r.owner_id,
			   r.name,
			   r.description,
			   r.is_fork,
			   r.created_at,
			   r.is_public
		FROM git_repositories AS r
				 JOIN users AS u ON r.owner_id = u.id
		WHERE u.login = $1
		  AND r.name = $2`,
		userLogin, repoName).Scan(
		&gitRepo.Id,
		&gitRepo.OwnerId,
		&gitRepo.Name,
		&gitRepo.Description,
		&gitRepo.IsFork,
		&gitRepo.CreatedAt,
		&gitRepo.IsPublic)

	switch {
	case err == sql.ErrNoRows:
		return gitRepo, entityerrors.DoesNotExist()
	case err != nil:
		return gitRepo, errors.Wrapf(err, "error while scanning in repository "+
			"for git repositories in GetByName with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	return gitRepo, nil
}

func (repo Repository) CheckReadAccess(currentUserId *int, userLogin, repoName string) (bool, error) {
	var haveAccess bool
	err := repo.db.QueryRow(`
		SELECT EXISTS(
		    SELECT * FROM users_git_repositories AS ugr 
		        JOIN git_repositories AS gr ON ugr.repository_id = gr.id
		    	JOIN users AS u ON gr.owner_id = u.id
		    WHERE gr.is_public = TRUE OR u.login = $1 AND gr.name = $2 AND ugr.user_id = $3 
		    )`, userLogin, repoName, currentUserId).Scan(&haveAccess)
	if err != nil {
		return false, errors.Wrapf(err, "error in repository for git repositories in CheckReadAccess "+
			"with currentUserId=%+v, userLogin=%s, repoName=%s", currentUserId, userLogin, repoName)
	}
	return haveAccess, nil
}
