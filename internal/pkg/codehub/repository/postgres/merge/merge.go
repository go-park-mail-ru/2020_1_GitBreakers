package merge

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type RepoPullReq struct {
	DB *sqlx.DB
}

func NewPullRequestRepository(db *sqlx.DB) RepoPullReq {
	return RepoPullReq{
		DB: db,
	}
}

func (repo RepoPullReq) CreateMR(request models.PullRequest) error {
	return nil
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
