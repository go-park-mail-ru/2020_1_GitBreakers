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
		"SELECT EXISTS(SELECT 1 FROM git_repository WHERE owner_id = $1 and name = $2)").Scan(&isRepoExists)
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
	_, err = tx.Exec(
		"INSERT INTO git_repository (owner_id, name, description, is_public, is_fork) VALUES ($1, $2, $3, $4, $5)",
		newRepo.OwnerId, newRepo.Name, newRepo.Description, newRepo.IsPublic, newRepo.IsFork)
	if err != nil {
		return errors.Wrap(err, "cannot create new git repository entity in database")
	}

	// Calculate path where git creates new repository on filesystem
	repoPath, err := createRepoPath(tx, newRepo.OwnerId, newRepo.Name)
	if err != nil {
		return err
	}

	// Create new bare repository aka 'git init --bare' on repoPath
	_, err = gogit.PlainInit(repo.reposDir+"/"+repoPath, true)
	if err == gogit.ErrRepositoryAlreadyExists {
		return entityerrors.AlreadyExist()
	}
	return nil
}

func (repo Repository) GetReposByUserLogin(userLogin string, offset, limit int) (gitRepos []git.Repository, err error) {
	rows, err := repo.db.Query(
		`SELECT r.id,
				   r.owner_id,
				   r.name,
				   r.description,
				   r.is_fork,
				   r.created_at,
				   r.is_public
			FROM git_repository AS r
					 JOIN users AS u ON r.owner_id = u.id
				OFFSET $1
			LIMIT $2`,
		offset, limit)

	if err != nil {
		return nil, errors.Wrapf(err, "error while performing query in repository "+
			"for git repositories with userLogin=%s, offset=%d, limit=$d", userLogin, offset, limit)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.Wrap(err, closeErr.Error())
		}
	}()

	gitRepos = make([]git.Repository, 0, limit)
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
				"while scanning userLogin=%s, offset=%d, limit=$d", userLogin, offset, limit)
		}
		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}
