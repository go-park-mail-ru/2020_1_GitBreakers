package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type Repository interface {
	AddStar(userID int64, repoID int64) error
	DelStar(userID int64, repoID int64) error
	GetStarredRepo(userID int64) ([]gitmodels.Repository, error)

	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID int64) error //закрывает вопрос(!не удаляет из бд)
	GetIssuesList(repoID int64) ([]models.Issue, error)
	CheckAccessIssue(userID, issueID int64) (perm.Permission, error) ///вернет "","read","update","close"
	CheckAccessRepo(userID, repoID int64) (perm.Permission, error)
	//обозначение по порядку(не может читать, может
	//просмотреть, может изменить(коллаборатор), может
	//закрыть(автор или админ репозитория)
	GetIssues(issueID int64) (models.Issue, error)
}
