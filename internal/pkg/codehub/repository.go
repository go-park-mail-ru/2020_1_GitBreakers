package codehub

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type Repository interface {
	AddStar(userID int, repoID int) error
	DelStar(userID int, repoID int) error
	GetStarredRepo(userID int) ([]gitmodels.Repository, error)

	CreateIssue(issue models.Issue) error
	UpdateIssue(issue models.Issue) error
	CloseIssue(issueID int) error //закрывает вопрос(!не удаляет из бд)
	GetIssuesList(repoID int) ([]models.Issue, error)
	CheckAccessIssue(userID, issueID int) (perm.Permission, error) ///вернет "","read","update","close"
	CheckAccessRepo(userID, repoID int) (perm.Permission, error)
	//обозначение по порядку(не может читать, может
	//просмотреть, может изменить(коллаборатор), может
	//закрыть(автор или админ репозитория)
	GetIssues(issueID int) (models.Issue, error)
}
