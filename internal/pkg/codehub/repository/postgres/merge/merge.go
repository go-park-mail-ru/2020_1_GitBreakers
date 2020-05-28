package merge

import (
	"bytes"
	"database/sql"
	"fmt"
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
	toRemoteName   = "to"
	fromRemoteName = "from"
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

func (repo RepoPullReq) CreateMR(request models.PullRequest) (pr models.PullRequest, err error) {
	if _, isExist, err := repo.isExistRepoAndBranch(*request.FromRepoID, request.BranchFrom); err != nil {
		return models.PullRequest{}, err
	} else if !isExist {
		return models.PullRequest{}, entityerrors.DoesNotExist()
	}

	if _, isExist, err := repo.isExistRepoAndBranch(request.ToRepoID, request.BranchTo); err != nil {
		return models.PullRequest{}, err
	} else if !isExist {
		return models.PullRequest{}, entityerrors.DoesNotExist()
	}

	tx, err := repo.db.Begin()
	if err != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}

	defer func() {
		err = finishTransaction(tx, err)
	}()

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
		return models.PullRequest{}, errors.WithStack(err)
	}

	if storageErr := repo.createEmptyMRStorage(request); storageErr != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}

	// zero in fetchDepth means fetch full
	if _, err = repo.fullUpdatePullRequest(tx, request, 0); err != nil {
		mrRepoPath := repo.getMrDirByID(request.ID)
		if _, existsErr := os.Stat(mrRepoPath); !os.IsNotExist(existsErr) {
			if removeErr := os.RemoveAll(mrRepoPath); removeErr != nil {
				err = errors.WithMessage(err, removeErr.Error())
			}
		}
		return models.PullRequest{}, errors.WithStack(err)
	}

	request, err = getMRByID(tx, request.ID)
	if err != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}

	return request, nil
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

func (repo RepoPullReq) ApproveMerge(mrID int64, approver models.User) (err error) {
	pr, err := repo.GetMRByID(mrID)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		return entityerrors.DoesNotExist()
	case err != nil:
		return errors.WithStack(err)
	}

	mrStatus, err := repo.fullUpdatePullRequest(repo.db, pr, 0) // zero means fetch full
	if err != nil {
		return errors.WithStack(err)
	}
	if mrStatus != codehub.MRStatusOK && mrStatus != codehub.MRStatusNoChanges {
		return entityerrors.Conflict()
	}

	err = repo.mergePullRequestAndPushInToRepo(pr, approver)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = repo.db.Exec(`
			UPDATE merge_requests
			SET is_closed   = TRUE,
				is_accepted = TRUE,
			    closer_user_id = $2,
			    status = $3
			WHERE id = $1
			RETURNING id`,
		mrID,
		approver.ID,
		codehub.MRStatusMerged,
	)
	if err != nil {
		return errors.WithStack(err)
	}

	mrRepoPath := repo.getMrDirByID(pr.ID)
	if _, existsErr := os.Stat(mrRepoPath); !os.IsNotExist(existsErr) {
		if removeErr := os.RemoveAll(mrRepoPath); removeErr != nil {
			err = removeErr
		}
	}

	return errors.WithStack(err)
}

func (repo RepoPullReq) GetAllMRForUser(userID int64, limit int64, offset int64) (pullRequests models.PullReqSet, err error) {
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
				   AND status IN ($2, $3, $4, $5)`,
		repoID,
		codehub.MRStatusOK,
		codehub.MRStatusConflict,
		codehub.MRStatusUpToDate,
		codehub.MRStatusNoChanges,
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			err = errors.WithMessage(err, closeErr.Error())
		}
	}()

	pullRequests = make(models.PullReqSet, 0)

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
		_, err := repo.fullUpdatePullRequest(repo.db, request, 0) // zero means fetch full
		if err != nil {
			errorsSlice = append(errorsSlice, err)
		}
	}

	return errorsSlice
}

func (repo RepoPullReq) GetMRDiffByID(mrID int64) (models.PullRequestDiff, error) {
	var diffModel models.PullRequestDiff

	err := repo.db.QueryRow(`
			SELECT diff, 
			       status
			FROM merge_requests
			WHERE id = $1`,
		mrID,
	).Scan(
		&diffModel.Diff,
		&diffModel.Status,
	)
	switch {
	case err == sql.ErrNoRows:
		return models.PullRequestDiff{}, entityerrors.DoesNotExist()
	case err != nil:
		return models.PullRequestDiff{}, errors.WithStack(err)
	}

	return diffModel, nil
}

func (repo RepoPullReq) GetMRByID(mrID int64) (models.PullRequest, error) {
	return getMRByID(repo.db, mrID)
}

func scanPullReq(rows *sql.Rows) (models.PullReqSet, error) {
	pullRequests := make(models.PullReqSet, 0)

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

func getMRByID(querier SQLInterfaces.Querier, mrID int64) (models.PullRequest, error) {
	rows, err := querier.Query(`
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

func (repo RepoPullReq) isExistRepoAndBranch(repoID int64, branchName string) (string, bool, error) {
	branchHash, err := repo.gitRepo.GetBranchHashIfExistInRepoByID(repoID, branchName)
	switch {
	case errors.Is(err, entityerrors.DoesNotExist()):
		return branchHash, false, nil
	case err != nil:
		return branchHash, false, errors.WithStack(err)
	}

	return branchHash, branchHash != "", nil
}

func (repo RepoPullReq) updateMRGitStorage(request models.PullRequest, fetchDepth int) error {
	mrRepoPath := repo.getMrDirByID(request.ID)

	gitRepo, err := git.OpenRepository(mrRepoPath)
	if err != nil {
		return err
	}

	if err := gitRepo.FetchBranchForce(toRemoteName, request.BranchTo, fetchDepth); err != nil {
		return err
	}

	if err := gitRepo.FetchBranchForce(fromRemoteName, request.BranchFrom, fetchDepth); err != nil {
		return err
	}

	toRemoteRefName := git.CreateRemoteRefName(toRemoteName, request.BranchTo)

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

// getGitDiff between baseRef and headRef
func (repo RepoPullReq) getGitDiff(request models.PullRequest, baseRef, headRef string,
	extraArgs ...string) (stdout, stderr string, err error) {

	mrRepoPath := repo.getMrDirByID(request.ID)

	gitDiffArgs := []string{
		"diff",
		"--full-index", // show full hash sum of commits
		"-M",           // detect renaming
		baseRef,
		headRef,
	}

	gitDiffArgs = append(gitDiffArgs, extraArgs...)

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

// getFormatPatch between firstRef and secondRef
func (repo RepoPullReq) getFormatPatch(request models.PullRequest, baseRef, headRef string,
	extraArgs ...string) (stdout, stderr string, err error) {

	mrRepoPath := repo.getMrDirByID(request.ID)

	gitFormatPatchArgs := []string{
		"format-patch",
		"--stdout", // print to stdout
		fmt.Sprintf("%s..%s", baseRef, headRef),
	}
	gitFormatPatchArgs = append(gitFormatPatchArgs, extraArgs...)

	outBuf, errBuf, err := process.ExecDir(-1, nil, mrRepoPath,
		fmt.Sprintf("calc new format patch for mr with path=%s", mrRepoPath),
		"git", gitFormatPatchArgs...,
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
	patch string, extraArgs ...string) (stdout, stderr string, err error) {

	mrRepoPath := repo.getMrDirByID(request.ID)

	gitApplyArgs := []string{
		"apply", "--check",
	}
	gitApplyArgs = append(gitApplyArgs, extraArgs...)

	inBuf := bytes.NewBufferString(patch)

	outBuf, errBuf, err := process.ExecDir(-1, inBuf, mrRepoPath,
		fmt.Sprintf("apply patch to repo with path=%s", mrRepoPath),
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
	toRemoteRefName := git.CreateRemoteRefName(toRemoteName, request.BranchTo)
	fromRemoteRefName := git.CreateRemoteRefName(fromRemoteName, request.BranchFrom)

	diff, _, err = repo.getGitDiff(request, toRemoteRefName, fromRemoteRefName)
	if err != nil {
		return "", err
	}

	err = updateMRDiffInDBByID(executer, request.ID, diff)
	if err != nil {
		return "", err
	}

	return diff, nil
}

func (repo RepoPullReq) checkMRConflicts(request models.PullRequest,
	toBranchHash, fromBranchHash string) (codehub.MergeRequestStatus, error) {

	if fromBranchHash == toBranchHash {
		return codehub.MRStatusUpToDate, nil
	}

	mrDir := repo.getMrDirByID(request.ID)

	gitRepo, err := git.OpenRepository(mrDir)
	if err != nil {
		return codehub.MRStatusNone, errors.WithStack(err)
	}

	mergeBases, err := gitRepo.MergeBase(fromBranchHash, toBranchHash)
	if err != nil {
		return codehub.MRStatusNone, errors.WithStack(err)
	}

	for i := range mergeBases {
		if fromBranchHash == mergeBases[i] {
			return codehub.MRStatusUpToDate, nil
		}
	}

	toRemoteRefName := git.CreateRemoteRefName(toRemoteName, request.BranchTo)
	fromRemoteRefName := git.CreateRemoteRefName(fromRemoteName, request.BranchFrom)

	patch, _, err := repo.getFormatPatch(request, toRemoteRefName, fromRemoteRefName)
	if err != nil {
		return codehub.MRStatusNone, errors.WithStack(err)
	}

	if patch == "" || patch == "\n" || patch == "\r\n" {
		return codehub.MRStatusNoChanges, nil
	}

	if _, _, err := repo.applyPatchCheck(request, patch); err != nil {
		return codehub.MRStatusConflict, nil
	}

	return codehub.MRStatusOK, nil
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

func (repo RepoPullReq) checkFromAndToBranches(executer SQLInterfaces.Executer,
	request models.PullRequest) (mrStatus codehub.MergeRequestStatus, toBranchHash, fromBranchHash string, err error) {

	var isExist bool

	toBranchHash, isExist, err = repo.isExistRepoAndBranch(request.ToRepoID, request.BranchTo)
	if err != nil {
		return codehub.MRStatusNone, "", "", errors.WithStack(err)
	}
	if !isExist {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID,
			codehub.MRStatusBadToBranch)

		return codehub.MRStatusBadToBranch, "", "", errors.WithStack(err)
	}

	fromBranchHash, isExist, err = repo.isExistRepoAndBranch(*request.FromRepoID, request.BranchFrom)
	if err != nil {
		return codehub.MRStatusNone, "", "", errors.WithStack(err)
	}
	if !isExist {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID,
			codehub.MRStatusBadFromBranch)

		return codehub.MRStatusBadFromBranch, "", "", errors.WithStack(err)
	}

	return codehub.MRStatusNone, toBranchHash, fromBranchHash, nil
}

func (repo RepoPullReq) fullUpdatePullRequest(executer SQLInterfaces.Executer,
	request models.PullRequest, fetchDepth int) (mrStatus codehub.MergeRequestStatus, err error) {

	var toBranchHash, fromBranchHash string

	// Check To And From repository branches
	if mrStatus, toBranchHash, fromBranchHash, err = repo.checkFromAndToBranches(executer, request); err != nil {
		return codehub.MRStatusNone, errors.WithStack(err)
	} else if mrStatus != codehub.MRStatusNone {
		return mrStatus, nil
	}

	// Fetching updated
	if updErr := repo.updateMRGitStorage(request, fetchDepth); updErr != nil {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID, codehub.MRStatusError)
		if err != nil {
			updErr = errors.WithMessage(updErr, err.Error())
		}
		return codehub.MRStatusError, errors.WithStack(updErr)
	}

	// Renew diff and store it in db
	_, renewDifErr := repo.renewMRDiff(executer, request)
	if renewDifErr != nil {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID, codehub.MRStatusError)
		if err != nil {
			renewDifErr = errors.WithMessage(renewDifErr, err.Error())
		}
		return codehub.MRStatusError, errors.WithStack(renewDifErr)
	}

	// Check mr and get new MR status
	if mrStatus, err = repo.checkMRConflicts(request, toBranchHash, fromBranchHash); err != nil {
		return codehub.MRStatusNone, errors.WithStack(err)
	}

	// Update status in db
	statusUpdateErr := updateMergeRequestsStatusById(executer, request.ID, mrStatus)
	if statusUpdateErr != nil {
		err := repo.forceCloseMRAndRemoveMRStorage(executer, request.ID, codehub.MRStatusError)
		if err != nil {
			statusUpdateErr = errors.WithMessage(statusUpdateErr, err.Error())
		}
		return codehub.MRStatusError, errors.WithStack(statusUpdateErr)
	}

	return mrStatus, nil
}

func (repo RepoPullReq) mergePullRequestAndPushInToRepo(pr models.PullRequest, merger models.User) (err error) {
	mrRepoPath := repo.getMrDirByID(pr.ID)

	fromRefName := git.CreateRemoteRefName(fromRemoteName, pr.BranchFrom)

	gitMergeArgs := []string{
		"merge",
		"--no-ff",
		"--no-commit",
		fromRefName,
	}

	var stderr []byte
	// Merge changes from head branch.
	if _, stderr, err = process.ExecDir(-1, nil, mrRepoPath,
		fmt.Sprintf("PullRequest.Merge (git merge --no-ff --no-commit): %s", mrRepoPath),
		"git", gitMergeArgs...); err != nil {

		var errOut string
		if stderr != nil {
			errOut = string(stderr)
		}
		return errors.WithStack(fmt.Errorf("git merge --no-ff --no-commit [%s]: %v - %s", mrRepoPath, err, errOut))
	}

	gitCommitArgs := []string{
		"commit",
		fmt.Sprintf("--author='%s <%s>'", merger.Name, merger.Email),
		"-m",
		fmt.Sprintf("Merge branch '%s' of %s/%s into %s",
			pr.BranchFrom, *pr.FromAuthorLogin, *pr.FromRepoName, pr.BranchTo),
		"-m",
		pr.Title,
	}
	// Commit merged changed
	if _, stderr, err = process.ExecDir(-1, nil, mrRepoPath,
		fmt.Sprintf("PullRequest.Merge (git merge): %s", mrRepoPath),
		"git", gitCommitArgs...); err != nil {

		var errOut string
		if stderr != nil {
			errOut = string(stderr)
		}
		return errors.WithStack(fmt.Errorf("git commit [%s]: %v - %s", mrRepoPath, err, errOut))
	}

	gitPushArgs := []string{
		"push",
		toRemoteName,
		fmt.Sprintf("HEAD:%s", pr.BranchTo),
	}
	// Push this changes to target repo
	if _, stderr, err = process.ExecDir(-1, nil, mrRepoPath,
		fmt.Sprintf("PullRequest.Merge (git push): %s", mrRepoPath),
		"git", gitPushArgs...); err != nil {

		var errOut string
		if stderr != nil {
			errOut = string(stderr)
		}
		return errors.WithStack(fmt.Errorf("git push: %s", errOut))
	}

	return nil
}
