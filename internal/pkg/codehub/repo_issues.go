package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type RepoIssueI interface {
	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID int64) error //закрывает вопрос(!не удаляет из бд)
	GetIssuesList(repoID int64, limit int64, offset int64) ([]models.Issue, error)
	CheckAccessIssue(userID, issueID int64) (perm.Permission, error) ///вернет "","read","update","close"
	CheckAccessRepo(userID, repoID int64) (perm.Permission, error)
	GetIssue(issueID int64) (models.Issue, error)
}
