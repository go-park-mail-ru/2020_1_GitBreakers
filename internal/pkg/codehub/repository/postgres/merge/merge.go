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

func (repo RepoPullReq) CreatePullReq(request models.PullRequest) error {
	return nil
}

func (repo RepoPullReq) GetAllPullReqOut(repoID int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) GetAllPullReqIn(repoID int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) ApproveMerge(pullReqID int64) error {
	return nil
}

func (repo RepoPullReq) GetOpenedPullReqForUser(userID int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) RejectPullReq(mrID int64) error {
	return nil
}
