package merge

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoPullReq struct {
}

func (repo RepoPullReq) CreatePullReq(authorID int64, hash1 string, fromID int64, hash2 string, toID int64) error {
	return nil
}

func (repo RepoPullReq) GetAllPullReq(repoID int64) (models.PullReqSet, error) {
	return models.PullReqSet{}, nil
}

func (repo RepoPullReq) ApproveMerge(pullReqID int64, userID int64) error {
	return nil
}

func (repo RepoPullReq) GetOnePullReq(pullReqID int64, userID int64) error {
	return nil
}
