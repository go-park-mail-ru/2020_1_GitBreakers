package merge

import (
	"bytes"
	"database/sql"
	"fmt"
	gogitPlumbing "github.com/go-git/go-git/v5/plumbing"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	gitPackage "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/process"
	SQLInterfaces "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"os"
	"path"
	"strconv"
)

const (
	// TODO use user service for this
	checkAuthorExistSQL = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	toRemoteName        = "to"
	fromRemoteName      = "from"
)

type RepoPullReq struct {
	db         *sqlx.DB
	gitRepo    gitPackage.GitRepoI
	pullReqDir string
}

func NewPullRequestRepository(db *sqlx.DB, gitRepo gitPackage.GitRepoI, pullReqDir string) RepoPullReq {
	return RepoPullReq{
		db:         db,
		gitRepo:    gitRepo,
		pullReqDir: pullReqDir,
	}
}

func (repo RepoPullReq) CreateMR(request models.PullRequest) (err error) {
	var isExist bool

	if request.FromRepoID == nil ||
		*request.FromRepoID == request.ToRepoID && request.BranchFrom == request.BranchTo {
		return entityerrors.Invalid()
	}

	if isExist, err := repo.isExistRepoAndBranch(*request.FromRepoID, request.BranchFrom); err != nil {
		return err
	} else if !isExist {
		return entityerrors.DoesNotExist()
	}

	if isExist, err := repo.isExistRepoAndBranch(request.ToRepoID, request.BranchTo); err != nil {
		return err
	} else if !isExist {
		return entityerrors.DoesNotExist()
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() {
		err = finishTransaction(tx, err)
	}()

	if err = tx.QueryRow(checkAuthorExistSQL, request.AuthorId).Scan(&isExist); err != nil {
		return err
	} else if !isExist {
		return entityerrors.DoesNotExist()
	}

	err = tx.QueryRow(`
			INSERT INTO merge_requests
			(author_id,
			 from_repository_id,
			 to_repository_id,
			 from_repository_branch,
			 to_repository_branch,
			 title,
			 message,
			 label)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		request.AuthorId,
		request.FromRepoID,
		request.ToRepoID,
		request.BranchFrom,
		request.BranchTo,
		request.Title,
		request.Message,
		request.Label,
	).Scan(
		&request.ID,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	if storageErr := repo.createEmptyMRStorage(request); storageErr != nil {
		return errors.WithStack(err)
	}

	if err = repo.fullUpdatePullRequest(tx, request); err != nil {
		mrRepoPath := repo.getMrDirByID(request.ID)
		if _, existsErr := os.Stat(mrRepoPath); !os.IsNotExist(existsErr) {
			if removeErr := os.RemoveAll(mrRepoPath); removeErr != nil {
				err = errors.WithMessage(err, removeErr.Error())
			}
		}
		return err
	}

	return nil
}

func (repo RepoPullReq) GetAllMROut(repoID int64, limit int64, offset int64) (pullRequests models.PullReqSet, err error) {
	rows, err := repo.db.Query(`
				SELECT mrv.id,
					   mrv.author_id,
				       mrv.closer_user_id,
					   mrv.from_repository_id,
					   mrv.to_repository_id,
					   mrv.from_repository_branch,
					   mrv.to_repository_branch,
					   mrv.title,
					   mrv.message,
					   mrv.label,
				       mrv.status,
					   mrv.is_closed,
					   mrv.is_accepted,
					   mrv.created_at,
				       mrv.git_repository_from_name,
				       mrv.git_repository_to_name,
				       upvfrom.login,
				       upvto.login
				FROM merge_requests_view AS mrv 
					JOIN user_profile_view AS upvfrom ON mrv.git_repository_from_owner_id = upvfrom.id
					JOIN user_profile_view AS upvto ON mrv.git_repository_to_owner_id = upvto.id
				WHERE mrv.from_repository_id = $1 LIMIT $2
				OFFSET $3`,
		repoID, limit, offset,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.WithMessage(err, closeErr.Error())
		}
	}()

	if pullRequests, err = scanPullReq(rows); err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (repo RepoPullReq) GetAllMRIn(repoID int64, limit int64, offset int64) (pullRequests models.PullReqSet, err error) {
	rows, err := repo.db.Query(`
				SELECT mrv.id,
					   mrv.author_id,
				       mrv.closer_user_id,
					   mrv.from_repository_id,
					   mrv.to_repository_id,
					   mrv.from_repository_branch,
					   mrv.to_repository_branch,
					   mrv.title,
					   mrv.message,
					   mrv.label,
				       mrv.status,
					   mrv.is_closed,
					   mrv.is_accepted,
					   mrv.created_at,
				       mrv.git_repository_from_name,
				       mrv.git_repository_to_name,
				       upvfrom.login,
				       upvto.login
				FROM merge_requests_view AS mrv 
					JOIN user_profile_view AS upvfrom ON mrv.git_repository_from_owner_id = upvfrom.id
					JOIN user_profile_view AS upvto ON mrv.git_repository_to_owner_id = upvto.id
				WHERE mrv.to_repository_id = $1 LIMIT $2
				OFFSET $3`,
		repoID, limit, offset,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.WithMessage(err, closeErr.Error())
		}
	}()

	if pullRequests, err = scanPullReq(rows); err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (repo RepoPullReq) ApproveMerge(approverID, pullReqID int64) error {
	err := repo.db.QueryRow(`
			UPDATE merge_requests
			SET is_closed   = TRUE,
				is_accepted = TRUE,
			    closer_user_id = $2,
			    status = $3
			WHERE id = $1
			RETURNING id`,
		pullReqID,
		approverID,
		codehub.MRStatusMerged,
	).Scan(
		&pullReqID,
	)
	switch {
	case err == sql.ErrNoRows:
		return entityerrors.DoesNotExist()
	case err != nil:
		return errors.WithStack(err)
	}

	return nil
}

func (repo RepoPullReq) GetAllMRForUser(userID int64, limit int64, offset int64) (pullRequests models.PullReqSet, err error) {
	pullRequests = models.PullReqSet{}

	var isUserExist bool
	if err := repo.db.QueryRow(checkAuthorExistSQL, userID).Scan(&isUserExist); err != nil {
		return nil, err
	} else if !isUserExist {
		return nil, entityerrors.DoesNotExist()
	}

	rows, err := repo.db.Query(`
				SELECT mrv.id,
					   mrv.author_id,
				       mrv.closer_user_id,
					   mrv.from_repository_id,
					   mrv.to_repository_id,
					   mrv.from_repository_branch,
					   mrv.to_repository_branch,
					   mrv.title,
					   mrv.message,
					   mrv.label,
				       mrv.status,
					   mrv.is_closed,
					   mrv.is_accepted,
					   mrv.created_at,
				       mrv.git_repository_from_name,
				       mrv.git_repository_to_name,
				       upvfrom.login,
				       upvto.login
				FROM merge_requests_view AS mrv 
					JOIN user_profile_view AS upvfrom ON mrv.git_repository_from_owner_id = upvfrom.id
					JOIN user_profile_view AS upvto ON mrv.git_repository_to_owner_id = upvto.id
				WHERE mrv.author_id = $1 LIMIT $2
				OFFSET $3`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.WithMessage(err, closeErr.Error())
		}
	}()

	if pullRequests, err = scanPullReq(rows); err != nil {
		return nil, err
	}

	return pullRequests, nil
}

func (repo RepoPullReq) RejectMR(rejecterID, mrID int64) error {
	err := repo.db.QueryRow(`
			UPDATE merge_requests
			SET is_closed   = TRUE,
				is_accepted = FALSE,
			    closer_user_id = $2,
			    status = $3
			WHERE id = $1
			RETURNING id`,
		mrID,
		rejecterID,
		codehub.MRStatusRejected,
	).Scan(
		&mrID,
	)
	switch {
	case err == sql.ErrNoRows:
		return entityerrors.DoesNotExist()
	case err != nil:
		return errors.WithStack(err)
	}

	mrRepoPath := repo.getMrDirByID(mrID)

	if _, existsErr := os.Stat(mrRepoPath); !os.IsNotExist(existsErr) {
		if removeErr := os.RemoveAll(mrRepoPath); removeErr != nil {
			err = removeErr
		}
	}

	return err
}

func (repo RepoPullReq) GetOpenedMRAssociatedWithRepoByRepoID(repoID int64) (pullRequests models.PullReqSet, err error) {
	rows, err := repo.db.Query(`
				SELECT mr.id,
					   mr.author_id,
				       mr.closer_user_id,
					   mr.from_repository_id,
					   mr.to_repository_id,
					   mr.from_repository_branch,
					   mr.to_repository_branch,
					   mr.title,
					   mr.message,
					   mr.label,
				       mr.status,
					   mr.is_closed,
					   mr.is_accepted,
					   mr.created_at
				FROM merge_requests AS mr
				WHERE (to_repository_id = $1
				   OR from_repository_id = $1)
				   AND status IN ($2, $3)`,
		repoID,
		codehub.MRStatusOK, codehub.MRStatusConflict,
	)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.WithMessage(err, closeErr.Error())
		}
	}()

	pullRequests = models.PullReqSet{}

	for rows.Next() {
		var pr models.PullRequest

		err := rows.Scan(
			&pr.ID,
			&pr.AuthorId,
			&pr.CloserUserId,
			&pr.FromRepoID,
			&pr.ToRepoID,
			&pr.BranchFrom,
			&pr.BranchTo,
			&pr.Title,
			&pr.Message,
			&pr.Label,
			&pr.Status,
			&pr.IsClosed,
			&pr.IsAccepted,
			&pr.CreatedAt,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		pullRequests = append(pullRequests, pr)
	}

	return pullRequests, nil
}

func (repo RepoPullReq) UpdateOpenedMRAssociatedWithRepoByRepoID(repoID int64) []error {
	requests, err := repo.GetOpenedMRAssociatedWithRepoByRepoID(repoID)
	if err != nil {
		return []error{err}
	}

	var errorsSlice []error
	for _, request := range requests {
		err := repo.fullUpdatePullRequest(repo.db, request)
		if err != nil {
			errorsSlice = append(errorsSlice, err)
		}
	}

	return errorsSlice
}

func (repo RepoPullReq) GetMRDiffByID(mrID int64) (string, error) {
	var diff []byte

	err := repo.db.QueryRow(`
			SELECT diff FROM merge_requests
			WHERE id = $1`,
		mrID,
	).Scan(
		&diff,
	)
	switch {
	case err == sql.ErrNoRows:
		return "", entityerrors.DoesNotExist()
	case err != nil:
		return "", errors.WithStack(err)
	}

	return string(diff), nil
}

func (repo RepoPullReq) GetMRByID(mrID int64) (models.PullRequest, error) {
	rows, err := repo.db.Query(`
				SELECT mrv.id,
					   mrv.author_id,
				       mrv.closer_user_id,
					   mrv.from_repository_id,
					   mrv.to_repository_id,
					   mrv.from_repository_branch,
					   mrv.to_repository_branch,
					   mrv.title,
					   mrv.message,
					   mrv.label,
				       mrv.status,
					   mrv.is_closed,
					   mrv.is_accepted,
					   mrv.created_at,
				       mrv.git_repository_from_name,
				       mrv.git_repository_to_name,
				       upvfrom.login,
				       upvto.login
				FROM merge_requests_view AS mrv 
					JOIN user_profile_view AS upvfrom ON mrv.git_repository_from_owner_id = upvfrom.id
					JOIN user_profile_view AS upvto ON mrv.git_repository_to_owner_id = upvto.id
				WHERE mrv.id = $1`,
		mrID,
	)
	if err != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.WithMessage(err, closeErr.Error())
		}
	}()

	var pullRequests models.PullReqSet
	if pullRequests, err = scanPullReq(rows); err != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}

	if len(pullRequests) == 0 {
		return models.PullRequest{}, entityerrors.DoesNotExist()
	}

	return pullRequests[0], nil
}

func scanPullReq(rows *sql.Rows) (models.PullReqSet, error) {
	var pullRequests models.PullReqSet

	for rows.Next() {
		var pr models.PullRequest

		err := rows.Scan(
			&pr.ID,
			&pr.AuthorId,
			&pr.CloserUserId,
			&pr.FromRepoID,
			&pr.ToRepoID,
			&pr.BranchFrom,
			&pr.BranchTo,
			&pr.Title,
			&pr.Message,
			&pr.Label,
			&pr.Status,
			&pr.IsClosed,
			&pr.IsAccepted,
			&pr.CreatedAt,
			&pr.FromRepoName,
			&pr.ToRepoName,
			&pr.FromAuthorLogin,
			&pr.ToAuthorLogin,
		)
		if err != nil {
			return nil, errors.WithStack(err)
		}

		pullRequests = append(pullRequests, pr)
	}

	return pullRequests, nil
}

func updateMergeRequestsStatusById(executer SQLInterfaces.Executer, mrID int64,
	status codehub.MergeRequestStatus) error {
	_, err := executer.Exec(`
			UPDATE merge_requests
			SET status = $2
			WHERE id = $1`,
		mrID,
		status,
	)
	return errors.WithStack(err)
}

func (repo RepoPullReq) getMrDirByID(prID int64) string {
	return path.Join(repo.pullReqDir, strconv.FormatInt(prID, 10))
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

func (repo RepoPullReq) isExistRepoAndBranch(repoID int64, branchName string) (bool, error) {
	isExist, err := repo.gitRepo.IsBranchExistInRepoByID(repoID, branchName)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		return false, nil
	case err != nil:
		return false, errors.WithStack(err)
	}

	return isExist, nil
}

func (repo RepoPullReq) updateMRGitStorage(request models.PullRequest) error {
	mrRepoPath := repo.getMrDirByID(request.ID)

	gitRepo, err := git.OpenRepository(mrRepoPath)
	if err != nil {
		return err
	}

	if err := gitRepo.FetchBranchForce(toRemoteName, request.BranchTo, 1); err != nil {
		return err
	}

	if err := gitRepo.FetchBranchForce(fromRemoteName, request.BranchFrom, 1); err != nil {
		return err
	}

	toRemoteRefName := gogitPlumbing.NewRemoteReferenceName(toRemoteName, request.BranchTo)

	if err := gitRepo.Checkout(toRemoteRefName, true); err != nil {
		return err
	}

	return nil
}

func (repo RepoPullReq) createEmptyMRStorage(request models.PullRequest) error {
	toRepoPath, err := repo.gitRepo.GetRepoPathByID(request.ToRepoID)
	if err != nil {
		return err
	}

	fromRepoPath, err := repo.gitRepo.GetRepoPathByID(*request.FromRepoID)
	if err != nil {
		return err
	}

	mrRepoPath := repo.getMrDirByID(request.ID)

	gitRepo, err := git.Init(mrRepoPath, false)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if removeErr := os.RemoveAll(mrRepoPath); removeErr != nil {
				err = errors.WithMessage(err, removeErr.Error())
			}
		}
	}()

	if _, err := gitRepo.AddRemote(git.ActionFileProtocol, toRemoteName, toRepoPath); err != nil {
		return err
	}

	if _, err := gitRepo.AddRemote(git.ActionFileProtocol, fromRemoteName, fromRepoPath); err != nil {
		return err
	}

	return nil
}

func (repo RepoPullReq) getDiff(request models.PullRequest) (stdout, stderr string, err error) {
	mrRepoPath := repo.getMrDirByID(request.ID)

	gitDiffArgs := []string{
		"diff",
		"--full-index", // for binary data such as images, executables and other blobs
		gogitPlumbing.NewRemoteReferenceName(toRemoteName, request.BranchTo).String(),
		gogitPlumbing.NewRemoteReferenceName(fromRemoteName, request.BranchFrom).String(),
	}

	outBuf, errBuf, err := process.ExecDir(-1, nil, mrRepoPath,
		fmt.Sprintf("calc new diff for mr with path=%s", mrRepoPath),
		"git", gitDiffArgs...,
	)

	if outBuf != nil {
		stdout = string(outBuf)
	}
	if errBuf != nil {
		stderr = string(errBuf)
	}

	return stdout, stderr, err
}

func (repo RepoPullReq) applyPatchCheck(request models.PullRequest,
	diff string) (stdout, stderr string, err error) {

	mrRepoPath := repo.getMrDirByID(request.ID)

	gitApplyArgs := []string{
		"apply", "--check",
	}

	inBuf := bytes.NewBufferString(diff)

	outBuf, errBuf, err := process.ExecDir(-1, inBuf, mrRepoPath,
		fmt.Sprintf("calc new diff for mr with path=%s", mrRepoPath),
		"git", gitApplyArgs...,
	)

	if outBuf != nil {
		stdout = string(outBuf)
	}
	if errBuf != nil {
		stderr = string(errBuf)
	}

	return stdout, stderr, err
}

func updateMRDiffInDBByID(executer SQLInterfaces.Executer, mrID int64, diff string) error {
	_, err := executer.Exec(`
			UPDATE merge_requests
			SET diff = $2
			WHERE id = $1`,
		mrID,
		[]byte(diff),
	)
	return err
}

func (repo RepoPullReq) renewMRDiff(executer SQLInterfaces.Executer, request models.PullRequest) (diff string, err error) {
	diff, _, err = repo.getDiff(request)
	if err != nil {
		return "", err
	}

	err = updateMRDiffInDBByID(executer, request.ID, diff)
	if err != nil {
		return "", err
	}

	return diff, nil
}

func (repo RepoPullReq) forceCloseMRAndRemoveMRStorage(executer SQLInterfaces.Executer, mrID int64,
	status codehub.MergeRequestStatus) error {

	_, err := executer.Exec(`
			UPDATE merge_requests
			SET is_closed   = TRUE,
				is_accepted = FALSE,
			    status = $2
			WHERE id = $1`,
		mrID,
		status,
	)

	if err != nil {
		return err
	}

	mrRepoPath := repo.getMrDirByID(mrID)

	if _, existsErr := os.Stat(mrRepoPath); !os.IsNotExist(existsErr) {
		if removeErr := os.RemoveAll(mrRepoPath); removeErr != nil {
			err = removeErr
		}
	}

	return err
}

func (repo RepoPullReq) fullUpdatePullRequest(executer SQLInterfaces.Executer, request models.PullRequest) error {
	// Check To repository branch
	isToExist, err := repo.gitRepo.IsBranchExistInRepoByID(request.ToRepoID, request.BranchTo)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isToExist {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID,
			codehub.MRStatusBadToBranch)
		if err != nil {
			return err
		}
	}

	// Check From repository branch
	isFromExist, err := repo.gitRepo.IsBranchExistInRepoByID(*request.FromRepoID, request.BranchFrom)
	if err != nil {
		return errors.WithStack(err)
	}
	if !isFromExist {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID,
			codehub.MRStatusBadFromBranch)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	// Fetching updated
	if updErr := repo.updateMRGitStorage(request); updErr != nil {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID, codehub.MRStatusError)
		if err != nil {
			updErr = errors.WithMessage(updErr, err.Error())
		}
		return errors.WithStack(updErr)
	}

	// Renew diff and store it in db
	diff, renewDifErr := repo.renewMRDiff(executer, request)
	if renewDifErr != nil {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID, codehub.MRStatusError)
		if err != nil {
			renewDifErr = errors.WithMessage(renewDifErr, err.Error())
		}
		return errors.WithStack(renewDifErr)
	}

	// Apply patch, if it conflicts, set status conflict
	mrStatus := codehub.MRStatusOK
	if _, _, patchErr := repo.applyPatchCheck(request, diff); patchErr != nil {
		mrStatus = codehub.MRStatusConflict
	}

	// Update status in db
	statusUpdateErr := updateMergeRequestsStatusById(executer, request.ID, mrStatus)
	if statusUpdateErr != nil {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID, codehub.MRStatusError)
		if err != nil {
			statusUpdateErr = errors.WithMessage(statusUpdateErr, err.Error())
		}
		return errors.WithStack(statusUpdateErr)
	}

	return nil
}
