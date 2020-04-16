package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
)

type Repository interface {
	AddStar(userID int, repoID int) error
	DelStar(userID int, repoID int) error
	GetStarredRepo(userID int) ([]gitmodels.Repository, error)

	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID int) error //закрывает вопрос(!не удаляет из бд)
	GetIssuesList(repoID int) ([]models.Issue, error)
}
