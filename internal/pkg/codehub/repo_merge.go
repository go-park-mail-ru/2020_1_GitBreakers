package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoMergeI interface {
	CreatePullReq(authorID int64, hash1 string, fromID int64, hash2 string, toID int64) error
	GetAllPullReq(repoID int64) (models.PullReqSet, error)
	ApproveMerge(pullReqID int64, userID int64) error
	GetOnePullReq(pullReqID int64, userID int64) error
}
