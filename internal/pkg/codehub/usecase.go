package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
)

type UCCodeHub interface {
	ModifyStar(star models.Star) error
	GetStarredRepo(userID int64) (models.RepoSet, error)
	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID int64, userID int64) error
	GetIssuesList(repoID int64, userID int64) (models.IssuesSet, error)
	GetIssue(issueID int64, userID int64) (models.Issue, error)
	GetNews(repoID int64, userID int64) (models.News, error)
}
