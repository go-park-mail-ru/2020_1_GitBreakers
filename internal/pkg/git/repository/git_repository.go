package repository

import (
	"database/sql"
	gogit "github.com/go-git/go-git/v5"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	gogitPlumbingObj "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	gitRepoPostfix = ".git"
)

type Repository struct {
	db       *sqlx.DB
	reposDir string
}

type queryer interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

func NewRepository(db *sqlx.DB, reposDir string) Repository {
	return Repository{db: db, reposDir: reposDir}
}

func convertToGitCommitModel(gogitCommit *gogitPlumbingObj.Commit) git.Commit {
	gogitCommitParentsHashes := make([]string, 0, len(gogitCommit.ParentHashes))
	for _, hash := range gogitCommit.ParentHashes {
		gogitCommitParentsHashes = append(gogitCommitParentsHashes, hash.String())
	}
	return git.Commit{
		CommitHash:        gogitCommit.Hash.String(),
		CommitAuthorName:  gogitCommit.Author.Name,
		CommitAuthorEmail: gogitCommit.Author.Email,
		CommitAuthorWhen:  gogitCommit.Author.When,
		CommitterName:     gogitCommit.Committer.Name,
		CommitterEmail:    gogitCommit.Committer.Email,
		CommitterWhen:     gogitCommit.Committer.When,
		TreeHash:          gogitCommit.TreeHash.String(),
		CommitParents:     gogitCommitParentsHashes,
	}
}

func getMimeTypeOfFile(gogitFile *gogitPlumbingObj.File) (string, error) {
	buffer := make([]byte, 512)

	entryFileReader, err := gogitFile.Blob.Reader()
	if err != nil {
		return "", errors.WithStack(err)
	}

	_, err = entryFileReader.Read(buffer)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return http.DetectContentType(buffer), nil
}

func isRepoExistsInDb(queryer queryer, ownerId int64, repoName string) (bool, error) {
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

func (repo Repository) convertToRepoPath(userLogin, repoName string) string {
	return repo.reposDir + "/" + userLogin + "/" + repoName + gitRepoPostfix
}

func (repo Repository) createRepoPath(queryer queryer, ownerId int64, repoName string) (string, error) {
	if repoName == "" {
		return "", entityerrors.Invalid()
	}
	var userLogin string
	err := queryer.QueryRow("SELECT login FROM users	 WHERE id = $1",
		ownerId).Scan(&userLogin)
	if err != nil {
		return "", err
	}
	return repo.convertToRepoPath(userLogin, repoName), nil
}

func (repo Repository) Create(newRepo git.Repository) (id int64, err error) {
	// New repo will be created by this path
	var repoPath string
	// Cleanup on error
	defer func() {
		if err != nil && repoPath != "" {
			if removeErr := os.RemoveAll(repoPath); removeErr != nil {
				err = errors.Wrapf(err, removeErr.Error())
			}
		}
	}()

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

	isRepoExist, err := isRepoExistsInDb(tx, newRepo.OwnerID, newRepo.Name)
	if err != nil {
		return -1, errors.Wrap(err, "error in create repository while checking if repository is not exits")
	}
	if isRepoExist {
		return -1, entityerrors.AlreadyExist()
	}

	var newRepoId int64
	// Create new db entity of git_repository
	err = tx.QueryRow(
		`INSERT INTO git_repositories (owner_id, name, description, is_public, is_fork) 
				VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		newRepo.OwnerID, newRepo.Name, newRepo.Description, newRepo.IsPublic, newRepo.IsFork).Scan(&newRepoId)
	if err != nil {
		return -1, errors.Wrapf(err, "cannot create new git repository entity in postgres, newRepo=%+v",
			newRepo)
	}

	_, err = tx.Exec("INSERT INTO users_git_repositories (user_id, repository_id, role) VALUES ($1, $2,$3)",
		newRepo.OwnerID, newRepoId, permission_types.OwnerAccess())
	if err != nil {
		return -1, errors.Wrapf(err, "cannot create new git repository entity in postgres, newRepo=%+v", newRepo)
	}

	// Calculate path where git creates new repository on filesystem
	repoPath, err = repo.createRepoPath(tx, newRepo.OwnerID, newRepo.Name)
	if err != nil {
		return -1, errors.Wrapf(err, "cannot create new git repository entity in postgres, newRepo=%+v", newRepo)
	}

	// Create new bare repository aka 'git init --bare' on repoPath
	_, err = gogit.PlainInit(repoPath, true)
	if err == gogit.ErrRepositoryAlreadyExists {
		return -1, entityerrors.AlreadyExist()
	}

	return newRepoId, nil
}

func (repo Repository) GetReposByUserLogin(requesterId *int64, userLogin string, offset, limit int64) ([]git.Repository, error) {
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
			"for git repositories in GetReposByUserLogin with userLogin=%s, requesterId=%+v offset=%d, limit=%d",
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
			&gitRepo.ID,
			&gitRepo.OwnerID,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.CreatedAt,
			&gitRepo.IsPublic)
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories "+
				"while scanning in GetReposByUserLogin userLogin=%s, offset=%d, limit=%d", userLogin, offset, limit)
		}
		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}

func (repo Repository) GetAnyReposByUserLogin(userLogin string, offset, limit int64) ([]git.Repository, error) {
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
			"for git repositories in GetAnyReposByUserLogin with userLogin=%s, offset=%d, limit=%d",
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
			&gitRepo.ID,
			&gitRepo.OwnerID,
			&gitRepo.Name,
			&gitRepo.Description,
			&gitRepo.IsFork,
			&gitRepo.CreatedAt,
			&gitRepo.IsPublic)
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories "+
				"while scanning in GetAnyReposByUserLogin userLogin=%s, offset=%d, limit=%d", userLogin, offset, limit)
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
		&gitRepo.ID,
		&gitRepo.OwnerID,
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

func (repo Repository) CheckReadAccess(currentUserId *int64, userLogin, repoName string) (bool, error) {
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
	gogitRepo, err := gogit.PlainOpen(repo.convertToRepoPath(userLogin, repoName))
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	gogitBranchesIterator, err := gogitRepo.Branches()
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	var gitRepoBranches []git.Branch

	err = gogitBranchesIterator.ForEach(func(reference *gogitPlumbing.Reference) error {
		gogitCommit, err := gogitPlumbingObj.GetCommit(gogitRepo.Storer, reference.Hash())
		if err != nil {
			return errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
				"with userLogin=%s, repoName=%s, branchName=%s",
				userLogin, repoName, reference.Name().String())
		}

		gitRepoBranches = append(gitRepoBranches,
			git.Branch{
				Name:   strings.TrimLeft(reference.Name().String(), gogitPlumbing.NewBranchReferenceName("").String()),
				Commit: convertToGitCommitModel(gogitCommit),
			},
		)

		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"after iterating by branches with userLogin=%s, repoName=%s", userLogin, repoName)
	}
	return gitRepoBranches, nil
}

func (repo Repository) GetByID(id int64) (git.Repository, error) {
	var gitRepo git.Repository
	err := repo.db.QueryRow(`
			SELECT id,
			       owner_id,
			       name,
			       description,
			       is_public,
			       is_fork, 
			       created_at
			FROM git_repositories WHERE id = $1
			`, id).Scan(
		&gitRepo.ID,
		&gitRepo.OwnerID,
		&gitRepo.Name,
		&gitRepo.Description,
		&gitRepo.IsPublic,
		&gitRepo.IsFork,
		&gitRepo.CreatedAt,
	)

	switch {
	case err == sql.ErrNoRows:
		return gitRepo, entityerrors.DoesNotExist()
	case err != nil:
		return gitRepo, errors.Wrapf(err, "error in repository for git repositories in GetByID "+
			"with id=%d", id)
	}

	return gitRepo, nil
}

func (repo Repository) FilesInCommitByPath(userLogin, repoName, commitHash, path string) ([]git.FileInCommit, error) {
	gogitRepo, err := gogit.PlainOpen(repo.convertToRepoPath(userLogin, repoName))
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	gogitCommit, err := gogitRepo.CommitObject(gogitPlumbing.NewHash(commitHash))
	switch {
	case err == gogitPlumbing.ErrObjectNotFound:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
			"with userLogin=%s, repoName=%s, commitHash=%s",
			userLogin, repoName, commitHash)
	}

	rootTree, err := gogitCommit.Tree()
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
			"while getting commit rootTree with userLogin=%s, repoName=%s, commitHash=%s",
			userLogin, repoName, commitHash)
	}

	var treeByPath *gogitPlumbingObj.Tree
	if path != "" && path != "/" && path != "./" { // TODO
		treeByPath, err = rootTree.Tree(path)
		if err != nil {
			return nil, entityerrors.DoesNotExist()
		}
	} else {
		treeByPath = rootTree
	}

	filesInCommit := make([]git.FileInCommit, 0, len(rootTree.Entries))
	for _, entry := range treeByPath.Entries {
		obj, err := gogitPlumbingObj.GetObject(gogitRepo.Storer, entry.Hash)
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
				"while getting object rootTree with userLogin=%s, repoName=%s, commitHash=%s, entryHash=%s",
				userLogin, repoName, commitHash, entry.Hash.String())
		}

		isBinary := true
		var fileSize int64 = -1
		var fileMIMEType string

		if obj.Type() == gogitPlumbing.BlobObject {
			entryFile, err := treeByPath.TreeEntryFile(&entry)
			if err != nil {
				return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
					"while getting tree entry file reader with userLogin=%s, repoName=%s, commitHash=%s, entry=%+v",
					userLogin, repoName, commitHash, entry)
			}

			isBinary, err = entryFile.IsBinary()
			if err != nil {
				return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
					"while checking IsBinary for file with userLogin=%s, repoName=%s, commitHash=%s, entry=%+v",
					userLogin, repoName, commitHash, entry)
			}

			fileSize = entryFile.Size
			fileMIMEType, err = getMimeTypeOfFile(entryFile)
			if err != nil {
				return nil, errors.Wrapf(err, "error in repository for git repositories in FilesInCommitByPath "+
					"while getting tree entry file reader with userLogin=%s, repoName=%s, commitHash=%s, entry=%+v",
					userLogin, repoName, commitHash, entry)
			}
		}

		fileInCommit := git.FileInCommit{
			Name:        entry.Name,
			FileType:    obj.Type().String(),
			FileMode:    git.FileMode(entry.Mode).String(),
			FileSize:    fileSize,
			IsBinary:    isBinary,
			ContentType: fileMIMEType,
			EntryHash:   entry.Hash.String(),
		}

		filesInCommit = append(filesInCommit, fileInCommit)
	}
	return filesInCommit, nil
}

func (repo Repository) GetCommitsByCommitHash(userLogin, repoName, commitHash string, offset, limit int64) ([]git.Commit, error) {
	gogitRepo, err := gogit.PlainOpen(repo.convertToRepoPath(userLogin, repoName))
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByCommitHash "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	gogitCommit, err := gogitRepo.CommitObject(gogitPlumbing.NewHash(commitHash))
	switch {
	case err == gogitPlumbing.ErrObjectNotFound:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByCommitHash "+
			"with userLogin=%s, repoName=%s, commitHash=%s",
			userLogin, repoName, commitHash)
	}

	commitIterator, err := gogitRepo.Log(
		&gogit.LogOptions{
			From:  gogitCommit.Hash,
			Order: gogit.LogOrderCommitterTime,
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByCommitHash "+
			"with userLogin=%s, repoName=%s, commitHash=%s",
			userLogin, repoName, commitHash)
	}
	defer commitIterator.Close()

	var gitCommits []git.Commit

	for limit > 0 {
		gogitCommit, err := commitIterator.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByCommitHash "+
				"with userLogin=%s, repoName=%s, commitHash=%s",
				userLogin, repoName, commitHash)
		}
		if offset > 0 {
			offset--
			continue
		}
		gitCommits = append(gitCommits, convertToGitCommitModel(gogitCommit))

		limit--
	}
	return gitCommits, nil
}

func (repo Repository) GetCommitsByBranchName(userLogin, repoName, branchName string, offset, limit int64) ([]git.Commit, error) {
	gogitRepo, err := gogit.PlainOpen(repo.convertToRepoPath(userLogin, repoName))
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByBranchName "+
			"with userLogin=%s, repoName=%s, branchName=%s", userLogin, repoName, branchName)
	}

	gogitBranch, err := gogitRepo.Reference(gogitPlumbing.NewBranchReferenceName(branchName), true)
	switch {
	case err == gogitPlumbing.ErrReferenceNotFound:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByBranchName "+
			"with userLogin=%s, repoName=%s, branchName=%s", userLogin, repoName, branchName)
	}

	gogitCommit, err := gogitRepo.CommitObject(gogitBranch.Hash())
	switch {
	case err == gogitPlumbing.ErrObjectNotFound:
		return nil, entityerrors.DoesNotExist()
	case err != nil:
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetCommitsByBranchName "+
			"with userLogin=%s, repoName=%s, branchName=%s, branchHash=%s",
			userLogin, repoName, branchName, gogitBranch.Hash().String())
	}

	return repo.GetCommitsByCommitHash(userLogin, repoName, gogitCommit.Hash.String(), offset, limit)
}

func (repo Repository) GetFileByPath(userLogin, repoName, commitHash, path string) (file git.FileCommitted, err error) {
	defer func() {
		if err != nil && err != entityerrors.DoesNotExist() && err != entityerrors.Invalid() {
			err = errors.Wrapf(err, "error in repository for git repositories in GetFileByPath "+
				"with userLogin=%s, repoName=%s, commitHash=%s, path=%s",
				userLogin, repoName, commitHash, path)
		}
	}()

	gogitRepo, err := gogit.PlainOpen(repo.convertToRepoPath(userLogin, repoName))
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return file, entityerrors.DoesNotExist()
	case err != nil:
		return file, errors.WithStack(err)
	}

	gogitCommit, err := gogitRepo.CommitObject(gogitPlumbing.NewHash(commitHash))
	switch {
	case err == gogitPlumbing.ErrObjectNotFound:
		return file, entityerrors.DoesNotExist()
	case err != nil:
		return file, err
	}

	rootTree, err := gogitCommit.Tree()
	if err != nil {
		return file, err
	}

	gogitFile, err := rootTree.File(path)
	switch {
	case err == gogitPlumbingObj.ErrFileNotFound:
		return file, entityerrors.DoesNotExist()
	case err != nil:
		return file, err
	}

	isBinary, err := gogitFile.IsBinary()
	if err != nil {
		return file, err
	}

	if gogitFile.Type() != gogitPlumbing.BlobObject || isBinary {
		return file, entityerrors.Invalid()
	}

	fileMIMEType, err := getMimeTypeOfFile(gogitFile)
	if err != nil {
		return file, err
	}

	gogitFileReader, err := gogitFile.Blob.Reader()
	if err != nil {
		return file, err
	}

	gogitFileContent, err := ioutil.ReadAll(gogitFileReader)
	if err != nil {
		return file, err
	}

	file = git.FileCommitted{
		FileInfo: git.FileInCommit{
			Name:        gogitFile.Name,
			FileType:    gogitFile.Type().String(),
			FileMode:    git.FileMode(gogitFile.Mode).String(),
			FileSize:    gogitFile.Size,
			IsBinary:    isBinary,
			ContentType: fileMIMEType,
			EntryHash:   gogitFile.Hash.String(),
		},
		Content: string(gogitFileContent),
	}

	return file, nil
}
