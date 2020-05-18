package merge

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
)

const (
	checkRepoExistSQL   = `SELECT EXISTS(SELECT 1 FROM git_repositories WHERE id = $1)`
	checkAuthorExistSQL = `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
)

type RepoPullReq struct {
	db      *sqlx.DB
	gitRepo git.GitRepoI
}

func NewPullRequestRepository(db *sqlx.DB) RepoPullReq {
	return RepoPullReq{
		db: db,
	}
}

func (repo RepoPullReq) CreateMR(request models.PullRequest) error {
	var isExist bool

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
	return err
}

func (repo RepoPullReq) GetAllMROut(repoID int64, limit int64, offset int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) GetAllMRIn(repoID int64, limit int64, offset int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) ApproveMerge(pullReqID int64) error {
	return nil
}

func (repo RepoPullReq) GetOpenedMRForUser(userID int64, limit int64, offset int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) RejectMR(mrID int64) error {
	return nil
}
