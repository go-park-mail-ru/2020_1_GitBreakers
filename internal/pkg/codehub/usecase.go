package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
)

type UCCodeHubI interface {
	ModifyStar(star models.Star) error
	GetStarredRepo(userID, limit, offset int64) (models.RepoSet, error)
	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID, userID int64) error
	GetIssuesList(repoID, userID, limit, offset int64) (models.IssuesSet, error)
	GetIssue(issueID, userID int64) (models.Issue, error)
	GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error)
}


