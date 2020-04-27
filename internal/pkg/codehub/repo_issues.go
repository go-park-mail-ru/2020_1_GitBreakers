package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type RepoIssueI interface {
	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	// CloseIssue return entityerrors.Invalid() if issueID is not valid ot this issue already closed
	CloseIssue(issueID int64) error
	GetIssuesList(repoID int64, limit int64, offset int64) ([]models.Issue, error)
	GetOpenedIssuesList(repoID int64, limit int64, offset int64) ([]models.Issue, error)
	GetClosedIssuesList(repoID int64, limit int64, offset int64) ([]models.Issue, error)
	CheckAccessIssue(userID, issueID int64) (perm.Permission, error) // TODO вернет "","read","update","close"
	CheckAccessRepo(userID, repoID int64) (perm.Permission, error) // TODO
	GetIssue(issueID int64) (models.Issue, error)
}
