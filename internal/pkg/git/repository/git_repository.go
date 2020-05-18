package repository

import (
	"database/sql"
	gogit "github.com/go-git/go-git/v5"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	gogitPlumbingObj "github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/jmoiron/sqlx"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

const (
	gitRepositoryNameSuffix = ".git"
)

type Repository struct {
	db       *sqlx.DB
	reposDir string
}

type queryer interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type execqueryer interface {
	execer
	queryer
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

	if err != nil && err != io.EOF {
		return "", errors.WithStack(err)
	}

	return http.DetectContentType(buffer), nil
}

func isRepoExistsInDbByOwnerId(queryer queryer, ownerId int64, repoName string) (bool, error) {
	if repoName == "" {
		return false, entityerrors.Invalid()
	}
	var isRepoExists bool
	err := queryer.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM git_repositories WHERE owner_id = $1 and name = $2)",
		ownerId,
		repoName,
	).Scan(
		&isRepoExists,
	)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return isRepoExists, nil
}

func isRepoExistsInDbByOwnerLogin(queryer queryer, ownerLogin string, repoName string) (bool, error) {
	if repoName == "" {
		return false, entityerrors.Invalid()
	}
	var isRepoExists bool
	err := queryer.QueryRow(`
			SELECT EXISTS(
						   SELECT 1
						   FROM git_repository_user_view
						   WHERE user_login = $1
							 AND name = $2
					   )`,
		ownerLogin,
		repoName,
	).Scan(
		&isRepoExists,
	)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return isRepoExists, nil
}

func getParentInfoById(queryer queryer, parentID *int64) (parentInfo git.ParentRepositoryInfo, err error) {
	if parentID == nil {
		return
	}

	err = queryer.QueryRow(`
			SELECT 	id,
					owner_id,
			       	name,
					user_login
			FROM git_repository_user_view WHERE id = $1`,
		parentID,
	).Scan(
		&parentInfo.ID,
		&parentInfo.OwnerID,
		&parentInfo.Name,
		&parentInfo.AuthorLogin,
	)
	if err != nil {
		return git.ParentRepositoryInfo{}, errors.WithStack(err)
	}

	return parentInfo, nil
}

func getByIdQ(queryer queryer, id int64) (git.Repository, error) {
	var gitRepo git.Repository
	err := queryer.QueryRow(`
			SELECT 	id,
			       	owner_id,
			       	name,
			       	description,
			       	is_fork, 
			       	is_public,
			       	stars,
					forks,
					merge_requests_open,
			       	created_at,
					user_login,
					parent_id,
					parent_owner_id,
					parent_name,
					parent_user_login
			FROM git_repository_parent_user_view WHERE id = $1`, id).Scan(
		&gitRepo.ID,
		&gitRepo.OwnerID,
		&gitRepo.Name,
		&gitRepo.Description,
		&gitRepo.IsFork,
		&gitRepo.IsPublic,
		&gitRepo.Stars,
		&gitRepo.Forks,
		&gitRepo.MergeRequestsOpen,
		&gitRepo.CreatedAt,
		&gitRepo.AuthorLogin,
		&gitRepo.ParentRepositoryInfo.ID,
		&gitRepo.ParentRepositoryInfo.OwnerID,
		&gitRepo.ParentRepositoryInfo.Name,
		&gitRepo.ParentRepositoryInfo.AuthorLogin,
	)

	switch {
	case err == sql.ErrNoRows:
		return gitRepo, entityerrors.DoesNotExist()
	case err != nil:
		return gitRepo, errors.Wrapf(err, "error in repository for git repositories in getByIdQ "+
			"with id=%d", id)
	}

	return gitRepo, nil
}

func createNewRepoSQLEntity(queryer queryer, newRepo git.Repository) (newRepoId int64, err error) {
	err = queryer.QueryRow(
		`INSERT INTO git_repositories (owner_id, name, description, is_public, is_fork, parent_id) 
				VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		newRepo.OwnerID,
		newRepo.Name,
		newRepo.Description,
		newRepo.IsPublic,
		newRepo.IsFork,
		newRepo.ParentRepositoryInfo.ID,
	).Scan(
		&newRepoId,
	)

	if err != nil {
		return 0, errors.WithStack(err)
	}

	return newRepoId, nil
}

func createNewPermissionSQLEntity(exec execer, userID, repoID int64,
	role permission_types.Permission) error {

	_, err := exec.Exec(
		`	INSERT INTO users_git_repositories (user_id, repository_id, role)
				VALUES ($1, $2, $3)`,
		userID,
		repoID,
		role,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func finishTransaction(tx *sql.Tx, err error) error {
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			err = errors.WithMessage(err, rollbackErr.Error())
		}
	} else if commitErr := tx.Commit(); commitErr != nil {
		err = errors.WithMessage(commitErr, commitErr.Error())
	}

	return err
}

func (repo Repository) convertToRepoPath(userLogin, repoName string) string {
	return repo.reposDir + "/" + userLogin + "/" + repoName + gitRepositoryNameSuffix
}

func (repo Repository) createRepoPath(queryer queryer, ownerId int64, repoName string) (string, error) {
	if repoName == "" {
		return "", entityerrors.Invalid()
	}
	var userLogin string
	err := queryer.QueryRow(
		"SELECT login FROM users WHERE id = $1",
		ownerId,
	).Scan(
		&userLogin,
	)
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
		return 0, errors.Wrap(err, "cannot begin transaction in create repository")
	}
	// Transaction cleanups
	defer func() {
		err = finishTransaction(tx, err)
	}()

	isRepoExist, err := isRepoExistsInDbByOwnerId(tx, newRepo.OwnerID, newRepo.Name)
	if err != nil {
		return 0, errors.Wrap(err, "error in create repository while checking if repository is not exits")
	}
	if isRepoExist {
		return -1, entityerrors.AlreadyExist()
	}

	// Create new db entity of git_repository
	newRepoId, err := createNewRepoSQLEntity(tx, newRepo)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot create new git repository entity in postgres, newRepo=%+v",
			newRepo)
	}

	err = createNewPermissionSQLEntity(tx, newRepo.OwnerID, newRepoId, perm.OwnerAccess())
	if err != nil {
		return 0, errors.Wrapf(err, "cannot create new git repository entity in postgres, newRepo=%+v", newRepo)
	}

	// Calculate path where git creates new repository on filesystem
	repoPath, err = repo.createRepoPath(tx, newRepo.OwnerID, newRepo.Name)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot create new git repository entity in postgres, newRepo=%+v", newRepo)
	}

	// Create new bare repository aka 'git init --bare' on repoPath
	_, err = gogit.PlainInit(repoPath, true)
	if err == gogit.ErrRepositoryAlreadyExists {
		return 0, entityerrors.AlreadyExist()
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
						repo.is_public,
						repo.stars,
						repo.forks,
        				repo.merge_requests_open,
						repo.created_at,
						repo.user_login,
						repo.parent_id,
						repo.parent_owner_id,
						repo.parent_name,
						repo.parent_user_login
				FROM git_repository_parent_user_view AS repo
					JOIN users_git_repositories AS ugr ON repo.id = ugr.repository_id
				WHERE repo.user_login = $1
				AND (ugr.user_id = $2 AND ugr.role <> $3 OR repo.is_public = TRUE) OFFSET $4
				LIMIT $5`,
		userLogin, requesterId, perm.NoAccess(), offset, limit)

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
			&gitRepo.IsPublic,
			&gitRepo.Stars,
			&gitRepo.Forks,
			&gitRepo.MergeRequestsOpen,
			&gitRepo.CreatedAt,
			&gitRepo.AuthorLogin,
			&gitRepo.ParentRepositoryInfo.ID,
			&gitRepo.ParentRepositoryInfo.OwnerID,
			&gitRepo.ParentRepositoryInfo.Name,
			&gitRepo.ParentRepositoryInfo.AuthorLogin,
		)
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
				repo.is_public,
		       	repo.stars,
		       	repo.forks,
		       	repo.merge_requests_open,
		       	repo.created_at,
		       	repo.user_login,
    			repo.parent_id,
				repo.parent_owner_id,
				repo.parent_name,
				repo.parent_user_login
		FROM git_repository_parent_user_view AS repo
			WHERE repo.user_login = $1 OFFSET $2
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
			&gitRepo.IsPublic,
			&gitRepo.Stars,
			&gitRepo.Forks,
			&gitRepo.MergeRequestsOpen,
			&gitRepo.CreatedAt,
			&gitRepo.AuthorLogin,
			&gitRepo.ParentRepositoryInfo.ID,
			&gitRepo.ParentRepositoryInfo.OwnerID,
			&gitRepo.ParentRepositoryInfo.Name,
			&gitRepo.ParentRepositoryInfo.AuthorLogin,
		)
		if err != nil {
			return nil, errors.Wrapf(err, "error in repository for git repositories "+
				"while scanning in GetAnyReposByUserLogin userLogin=%s, offset=%d, limit=%d",
				userLogin, offset, limit)
		}

		gitRepos = append(gitRepos, gitRepo)
	}
	return gitRepos, nil
}

func (repo Repository) GetByName(userLogin, repoName string) (git.Repository, error) {
	var gitRepo git.Repository

	err := repo.db.QueryRow(`
			SELECT 	repo.id,
					repo.owner_id,
					repo.name,
					repo.description,
					repo.is_fork,
					repo.is_public,
					repo.stars,
			       	repo.forks,
			       	repo.merge_requests_open,
					repo.created_at,
			        repo.user_login,
					repo.parent_id,
					repo.parent_owner_id,
					repo.parent_name,
					repo.parent_user_login
			FROM git_repository_parent_user_view AS repo
			WHERE repo.user_login = $1 AND repo.name = $2`,
		userLogin,
		repoName,
	).Scan(
		&gitRepo.ID,
		&gitRepo.OwnerID,
		&gitRepo.Name,
		&gitRepo.Description,
		&gitRepo.IsFork,
		&gitRepo.IsPublic,
		&gitRepo.Stars,
		&gitRepo.Forks,
		&gitRepo.MergeRequestsOpen,
		&gitRepo.CreatedAt,
		&gitRepo.AuthorLogin,
		&gitRepo.ParentRepositoryInfo.ID,
		&gitRepo.ParentRepositoryInfo.OwnerID,
		&gitRepo.ParentRepositoryInfo.Name,
		&gitRepo.ParentRepositoryInfo.AuthorLogin,
	)

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
	isExist, err := isRepoExistsInDbByOwnerLogin(repo.db, userLogin, repoName)
	if err != nil {
		return false, err
	}
	if !isExist {
		return false, entityerrors.DoesNotExist()
	}

	var haveAccess bool
	err = repo.db.QueryRow(`
		SELECT EXISTS(
				SELECT 1
				FROM users_git_repositories_view AS ugrv
						 JOIN user_profile_view AS own
							  ON ugrv.git_repository_owner_id = own.id
								  AND own.login = $1
								  AND ugrv.git_repository_name = $2
				WHERE ugrv.git_repository_is_public = TRUE
				   OR ugrv.user_id = $3
		    )`, userLogin, repoName, currentUserId).Scan(&haveAccess)
	if err != nil {
		return false, errors.Wrapf(err, "error in repository for git repositories in CheckReadAccess "+
			"with currentUserId=%+v, userLogin=%s, repoName=%s", currentUserId, userLogin, repoName)
	}
	return haveAccess, nil
}

func (repo Repository) GetPermission(currentUserId *int64, userLogin, repoName string) (perm.Permission, error) {
	var permissionRole string
	err := repo.db.QueryRow(`
		SELECT ugr.role FROM users_git_repositories AS ugr 
		        JOIN git_repositories AS gr ON ugr.repository_id = gr.id
		    	JOIN users AS u ON gr.owner_id = u.id
		    WHERE u.login = $1 AND gr.name = $2 AND ugr.user_id = $3`,
		userLogin, repoName, currentUserId).Scan(&permissionRole)

	switch {
	case err == sql.ErrNoRows:
		return perm.NoAccess(), nil
	case err != nil:
		return perm.NoAccess(), errors.Wrapf(err, "error in repository for git repositories in GetPermission "+
			"with userLogin=%s, repoName=%s", userLogin, repoName)
	}

	return perm.Permission(permissionRole), nil
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

	err = gogitBranchesIterator.ForEach(
		func(reference *gogitPlumbing.Reference) error {
			gogitCommit, err := gogitPlumbingObj.GetCommit(gogitRepo.Storer, reference.Hash())
			if err != nil {
				return errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
					"with userLogin=%s, repoName=%s, branchName=%s",
					userLogin, repoName, reference.Name().String())
			}

			referenceName := reference.Name().String()

			branchName := strings.TrimPrefix(referenceName, gogitPlumbing.NewBranchReferenceName("").String())

			gitRepoBranches = append(gitRepoBranches,
				git.Branch{
					Name:   branchName,
					Commit: convertToGitCommitModel(gogitCommit),
				},
			)

			return nil
		},
	)
	if err != nil {
		return nil, errors.Wrapf(err, "error in repository for git repositories in GetBranchesByName "+
			"after iterating by branches with userLogin=%s, repoName=%s", userLogin, repoName)
	}
	return gitRepoBranches, nil
}

func (repo Repository) GetByID(id int64) (git.Repository, error) {
	return getByIdQ(repo.db, id)
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
	if path != "" && path != "/" && path != "./" { // FIXME(nickeskov): this is a crunch to check root of repo
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

func (repo Repository) GetRepoHead(userLogin, repoName string) (defaultBranch git.Branch, err error) {
	gogitRepo, err := gogit.PlainOpen(repo.convertToRepoPath(userLogin, repoName))
	switch {
	case err == gogit.ErrRepositoryNotExists:
		return defaultBranch, entityerrors.DoesNotExist()
	case err != nil:
		return defaultBranch, errors.WithStack(err)
	}

	headReference, err := gogitRepo.Head()
	switch {
	case err == gogitPlumbing.ErrReferenceNotFound:
		return defaultBranch, entityerrors.ContentEmpty()
	case err != nil:
		return defaultBranch, errors.WithStack(err)
	}

	gogitCommit, err := gogitPlumbingObj.GetCommit(gogitRepo.Storer, headReference.Hash())
	switch {
	case err == gogitPlumbing.ErrObjectNotFound:
		return defaultBranch, entityerrors.DoesNotExist()
	case err != nil:
		return defaultBranch, errors.WithStack(err)
	}

	referenceName := headReference.Name().String()

	branchName := strings.TrimPrefix(referenceName, gogitPlumbing.NewBranchReferenceName("").String())

	defaultBranch = git.Branch{
		Name:   branchName,
		Commit: convertToGitCommitModel(gogitCommit),
	}

	return defaultBranch, nil
}

func (repo Repository) Fork(forkRepoName string, userID, repoBaseID int64) (err error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		err = finishTransaction(tx, err)
	}()

	forkFromRepo, err := getByIdQ(tx, repoBaseID)
	if err != nil {
		return err
	}
	if forkFromRepo.OwnerID == userID {
		return entityerrors.Conflict()
	}

	if isForkExist, err := isRepoExistsInDbByOwnerId(tx, userID, forkRepoName); err != nil {
		return err
	} else if isForkExist {
		return entityerrors.AlreadyExist()
	}

	newForkRepoPath, err := repo.createRepoPath(tx, userID, forkRepoName)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if _, existsErr := os.Stat(newForkRepoPath); !os.IsNotExist(existsErr) {
				if removeErr := os.RemoveAll(newForkRepoPath); removeErr != nil {
					err = errors.WithMessage(err, removeErr.Error())
				}
			}
		}
	}()

	newForkedRepo := git.Repository{
		OwnerID:     userID,
		Name:        forkRepoName,
		Description: forkFromRepo.Description,
		IsFork:      true,
		IsPublic:    forkFromRepo.IsPublic,
		ParentRepositoryInfo: git.ParentRepositoryInfo{
			ID: &forkFromRepo.ID,
		},
	}

	newForkedRepo.ID, err = createNewRepoSQLEntity(tx, newForkedRepo)
	if err != nil {
		return errors.WithStack(err)
	}

	err = createNewPermissionSQLEntity(tx, userID, newForkedRepo.ID, perm.OwnerAccess())
	if err != nil {
		return errors.WithStack(err)
	}

	if !newForkedRepo.IsPublic {
		_, err := tx.Exec(`
				INSERT INTO users_git_repositories (user_id, repository_id, role)
				SELECT user_id, $2, $3
				FROM users_git_repositories
				WHERE repository_id = $1`,
			forkFromRepo.ID,
			newForkedRepo.ID,
			perm.ReadAccess(),
		)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	forkFromRepoPath := repo.convertToRepoPath(forkFromRepo.AuthorLogin, forkFromRepo.Name)

	if err = copy.Copy(forkFromRepoPath, newForkRepoPath); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
