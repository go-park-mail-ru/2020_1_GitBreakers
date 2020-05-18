package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoMergeI interface {
	CreatePullReq(request models.PullRequest) error
	GetAllPullReqOut(repoID int64) (models.PullReqSet, error)
	GetAllPullReqIn(repoID int64) (models.PullReqSet, error)
	ApproveMerge(pullReqID int64) error
	GetOpenedPullReqForUser(userID int64) (models.PullReqSet, error)
	// GetOnePullReq(pullReqID int64) (models.PullRequest, error)
}
