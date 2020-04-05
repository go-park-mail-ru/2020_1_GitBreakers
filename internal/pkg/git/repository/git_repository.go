package repository

import (
	"database/sql"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	gogitPlumbingObj "github.com/go-git/go-git/v5/plumbing/object"
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

type queryer interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func NewRepository(db *sqlx.DB, reposDirs string) Repository {
	return Repository{db: db, reposDir: reposDirs}
}

func convertToRepoPath(userLogin, repoName string) string {
	return userLogin + gitPostfix + "/" + repoName
}

func createRepoPath(queryer queryer, ownerId int, repoName string) (string, error) {
	if repoName == "" {
		return "", entityerrors.Invalid()
	}
	var userLogin string
	err := queryer.QueryRow("SELECT login FROM users WHERE id = $1",
		ownerId).Scan(&userLogin)
	if err != nil {
		return "", err
	}
	return convertToRepoPath(userLogin, repoName), nil
}

func isRepoExistsInDb(queryer queryer, ownerId int, repoName string) (bool, error) {
	if repoName == "" {
		return false, entityerrors.Invalid()
	}
	var isRepoExists bool
	err := queryer.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM git_repositories WHERE owner_id = $1 and name = $2)",
		ownerId, repoName).Scan(&isRepoExists)
	if err != nil {
		return false, errors.Wrap(err, "error while checking if git repository exists")
	}
	return isRepoExists, nil
}

func (repo Repository) Create(newRepo git.Repository) (id int64, err error) {
	// Begin transaction
	tx, err := repo.db.Begin()
	if err != nil {
		return -1, errors.Wrap(err, "cannot begin transaction in create repository")
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
		return -1, errors.Wrap(err, "error in create repository while checking if repository is not exits")
	}
	if isRepoExist {
		return -1, entityerrors.AlreadyExist()
	}

	// Create new db entity of git_repository
	repoCreationResult, err := tx.Exec(
		`INSERT INTO git_repositories (owner_id, name, description, is_public, is_fork) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		newRepo.OwnerId, newRepo.Name, newRepo.Description, newRepo.IsPublic, newRepo.IsFork)
	if err != nil {
		return -1, errors.Wrapf(err, "cannot create new git repository entity in database, newRepo=%+v",
			newRepo)
	}

	newRepoId, err := repoCreationResult.LastInsertId()
	if err != nil {
		return -1, errors.Wrapf(err, "cannot get id from git repository entity from db (this entity " +
			"not exist in db now), newRepo=%+v", newRepo)
	}
	_, err = tx.Exec("INSERT INTO users_git_repositories (repository_id, user_id) VALUES ($1, $2)",
		newRepoId, newRepo.OwnerId)
	if err != nil {
		return -1, errors.Wrapf(err, "cannot create new git repository entity in database, newRepo=%+v", newRepo)
	}

	// Calculate path where git creates new repository on filesystem
	repoPath, err := createRepoPath(tx, newRepo.OwnerId, newRepo.Name)
	if err != nil {
		return -1, errors.Wrapf(err, "cannot create new git repository entity in database, newRepo=%+v", newRepo)
	}

	// Create new bare repository aka 'git init --bare' on repoPath
	_, err = gogit.PlainInit(repo.reposDir+"/"+repoPath, true)
	if err == gogit.ErrRepositoryAlreadyExists {
		return -1, entityerrors.AlreadyExist()
	}

	return newRepoId, nil
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

func (repo Repository) GetBranchesByName(userLogin, repoName string) ([]git.Branch, error) {
	gogitRepo, err := gogit.PlainOpen(userLogin + gitPostfix + "/" + repoName)
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	gogitBranchesIter, err := gogitRepo.Branches()
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	var gitRepoBranches []git.Branch

	err = gogitBranchesIter.ForEach(func(reference *gogitPlumbing.Reference) error {
		gogitCommit, err := gogitPlumbingObj.GetCommit(gogitRepo.Storer, reference.Hash())
		if err != nil {
			return errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
				"with userLogin=%s, repoName=%s, branchName=%s",
				userLogin, repoName, reference.Name().String())
		}

		gogitCommitParentHashes := make([]string, 0, len(gogitCommit.ParentHashes))
		for _, hash := range gogitCommit.ParentHashes {
			gogitCommitParentHashes = append(gogitCommitParentHashes, hash.String())
		}

		gitBranch := git.Branch{
			Name: reference.Name().String(),
			Commit: git.Commit{
				CommitHash:        gogitCommit.Hash.String(),
				CommitAuthorName:  gogitCommit.Author.Name,
				CommitAuthorEmail: gogitCommit.Author.Email,
				CommitAuthorWhen:  gogitCommit.Author.When,
				CommitterName:     gogitCommit.Committer.Name,
				CommitterEmail:    gogitCommit.Committer.Email,
				CommitterWhen:     gogitCommit.Committer.When,
				TreeHash:          gogitCommit.TreeHash.String(),
				CommitParents:     gogitCommitParentHashes,
			},
		}
		gitRepoBranches = append(gitRepoBranches, gitBranch)
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"after iterating by branches with userLogin=%s, repoName=%s", userLogin, repoName)
	}
	return gitRepoBranches, nil
}
