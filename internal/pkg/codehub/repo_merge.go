package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoMergeI interface {
	CreateMR(request models.PullRequest) error
	GetAllMROut(repoID int64, limit int64, offset int64) (models.PullReqSet, error)
	GetAllMRIn(repoID int64, limit int64, offset int64) (models.PullReqSet, error)
	ApproveMerge(approverID, pullReqID int64) error
	GetAllMRForUser(userID int64, limit int64, offset int64) (models.PullReqSet, error)
	RejectMR(rejecterID, mrID int64) error
	UpdateOpenedMergeRequestsStatusByRepoId(status MergeRequestStatus, repoID int64) error
	// GetOnePullReq(pullReqID int64) (models.PullRequest, error)
}

type MergeRequestStatus string

const (
	StatusOK          = "ok"
	StatusError       = "error"
	StatusMerged      = "merged"
	StatusRejected    = "rejected"
	StatusConflict    = "conflict"
	StatusNeedsUpdate = "needs_update"
)
