package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
)

type UCCodeHub interface {
	ModifyStar(star models.Star) error
	GetStarredRepo(userID int) (models.RepoSet, error)
	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID int, userID int) error
	GetIssuesList(repoID int, userID int) (models.IssuesSet, error)
	GetIssue(issueID int, userID int) (models.Issue, error)
}
