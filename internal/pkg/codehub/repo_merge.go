package codehub

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoMergeI interface {
	CreateMR(request models.PullRequest) (models.PullRequest, error)
	GetAllMROut(repoID int64, limit int64, offset int64) (models.PullReqSet, error)
	GetAllMRIn(repoID int64, limit int64, offset int64) (models.PullReqSet, error)
	GetMRDiffByID(mrID int64) (string, error)
	GetMRByID(mrID int64) (models.PullRequest, error)
	ApproveMerge(pullReqID int64, approver models.User) error
	GetAllMRForUser(userID int64, limit int64, offset int64) (models.PullReqSet, error)
	RejectMR(rejecterID, mrID int64) error
	GetOpenedMRAssociatedWithRepoByRepoID(repoID int64) (models.PullReqSet, error)
	UpdateOpenedMRAssociatedWithRepoByRepoID(repoID int64) []error
	// GetOnePullReq(pullReqID int64) (models.PullRequest, error)
}

type MergeRequestStatus string

const (
	MRStatusNone          MergeRequestStatus = ""                // This is invalid status, if this happens, something wrong
	MRStatusOK            MergeRequestStatus = "ok"              // Waiting for merge (can perform auto merge)
	MRStatusError         MergeRequestStatus = "error"           // MR closed by server error
	MRStatusMerged        MergeRequestStatus = "merged"          // MR closed by successful merge
	MRStatusRejected      MergeRequestStatus = "rejected"        // MR closed and changes not accepted
	MRStatusConflict      MergeRequestStatus = "conflict"        // Waiting for merge (conflict)
	MRStatusBadToBranch   MergeRequestStatus = "bad_to_branch"   // MR closed, because to branch not exist
	MRStatusBadFromBranch MergeRequestStatus = "bad_from_branch" // MR closed, because from branch not exist
)
