package merge

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	SQLInterfaces "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/sql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

const (
	checkRepoExistSQL   = `SELECT EXISTS(SELECT 1 FROM git_repositories WHERE id = $1)`
	checkAuthorExistSQL = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
)

type RepoPullReq struct {
	db         *sqlx.DB
	//gitRepo    git.GitRepoI
	pullReqDir string
}

func NewPullRequestRepository(db *sqlx.DB, pullReqDir string) RepoPullReq {
	return RepoPullReq{
		db:         db,
		pullReqDir: pullReqDir,
	}
}

func scanPullReq(rows *sql.Rows) (models.PullReqSet, error) {
	var pullRequests models.PullReqSet

	for rows.Next() {
		var pr models.PullRequest

		err := rows.Scan(
			&pr.ID,
			&pr.AuthorId,
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

func updateMergeRequestsStatusByRepoId(exec SQLInterfaces.Executer,
	status codehub.MergeRequestStatus, repoID int64) error {
	_, err := exec.Exec(`
			UPDATE merge_requests
			SET status = $1
			WHERE to_repository_id = $2
			   OR from_repository_id = $2`,
		status,
		repoID,
	)
	return err
}

func (repo RepoPullReq) CreateMR(request models.PullRequest) error {
	var isExist bool

	if request.FromRepoID == request.ToRepoID && request.BranchFrom == request.BranchTo {
		return entityerrors.Invalid()
	}

	if err := repo.db.QueryRow(checkRepoExistSQL, request.FromRepoID).Scan(&isExist); err != nil {
		return err
	} else if !isExist {
		return entityerrors.DoesNotExist()
	}

	if err := repo.db.QueryRow(checkRepoExistSQL, request.ToRepoID).Scan(&isExist); err != nil {
		return err
	} else if !isExist {
		return entityerrors.DoesNotExist()
	}

	if err := repo.db.QueryRow(checkAuthorExistSQL, request.AuthorId).Scan(&isExist); err != nil {
		return err
	} else if !isExist {
		return entityerrors.DoesNotExist()
	}

	_, err := repo.db.Exec(`
			INSERT INTO merge_requests
			(author_id,
			 from_repository_id,
			 to_repository_id,
			 from_repository_branch,
			 to_repository_branch,
			 title,
			 message,
			 label)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		request.AuthorId,
		request.FromRepoID,
		request.ToRepoID,
		request.BranchFrom,
		request.BranchTo,
		request.Title,
		request.Message,
		request.Label,
	)
	if err != nil {
		return errors.WithStack(err)
	}
	return err
}

func (repo RepoPullReq) GetAllMROut(repoID int64, limit int64, offset int64) (pullRequests models.PullReqSet, err error) {
	rows, err := repo.db.Query(`
				SELECT mrv.id,
					   mrv.author_id,
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

func (repo RepoPullReq) ApproveMerge(pullReqID int64) error {
	err := repo.db.QueryRow(`
			UPDATE merge_requests
			SET is_closed   = TRUE,
				is_accepted = TRUE
			WHERE id = $1
			RETURNING id`, pullReqID).Scan(&pullReqID)
	switch {
	case err == sql.ErrNoRows:
		return entityerrors.DoesNotExist()
	case err != nil:
		return errors.WithStack(err)
	}

	return nil
}

func (repo RepoPullReq) GetOpenedMRForUser(userID int64, limit int64, offset int64) (pullRequests models.PullReqSet, err error) {
	var isUserExist bool
	if err := repo.db.QueryRow(checkAuthorExistSQL, userID).Scan(&isUserExist); err != nil {
		return nil, err
	} else if !isUserExist {
		return nil, entityerrors.DoesNotExist()
	}

	rows, err := repo.db.Query(`
				SELECT mrv.id,
					   mrv.author_id,
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
				WHERE mrv.author_id = $1 AND mrv.is_closed = FALSE LIMIT $2
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

func (repo RepoPullReq) RejectMR(mrID int64) error {
	err := repo.db.QueryRow(`
			UPDATE merge_requests
			SET is_closed   = TRUE,
				is_accepted = FALSE
			WHERE id = $1
			RETURNING id`,
		mrID,
	).Scan(
		&mrID,
	)
	switch {
	case err == sql.ErrNoRows:
		return entityerrors.DoesNotExist()
	case err != nil:
		return errors.WithStack(err)
	}

	return nil
}

func (repo RepoPullReq) UpdateMergeRequestsStatusByRepoId(status codehub.MergeRequestStatus,
	repoID int64) error {
	return updateMergeRequestsStatusByRepoId(repo.db, status, repoID)
}
