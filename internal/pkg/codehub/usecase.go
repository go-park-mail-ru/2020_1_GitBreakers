package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type UCCodeHubI interface {
	ModifyStar(star models.Star) error
	GetStarredRepos(userID, limit, offset int64) (models.RepoSet, error)
	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID, userID int64) error
	GetIssuesList(repoID, userID, limit, offset int64) (models.IssuesSet, error)
	GetIssue(issueID, userID int64) (models.Issue, error)
	GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error)
	GetUserStaredList(repoID int64, limit int64, offset int64) (models.UserSet, error)
	Search(query, params string, limit, offset, userID int64) (interface{}, error)

	CreatePL(request models.PullRequest) error
	GetPLIn(repo gitmodels.Repository) (models.PullReqSet, error)
	GetPLOut(repo gitmodels.Repository) (models.PullReqSet, error)
	ApprovePL(plID int64, userID int64) error
	ClosePL(plID int64, userID int64) error
	GetAllMRUser(userID int64) (models.PullReqSet, error)
}
