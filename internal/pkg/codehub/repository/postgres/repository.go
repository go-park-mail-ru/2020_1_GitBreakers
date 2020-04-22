package postgres

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repo {
	return Repo{DB: db}
}
func (R *Repo) AddStar(userID int, repoID int) error {
	return nil
}

func (R *Repo) DelStar(userID int, repoID int) error {
	return nil
}

func (R *Repo) GetStarredRepo(userID int) ([]gitmodels.Repository, error) {
	return []gitmodels.Repository{}, nil
}

func (R *Repo) CreateIssue(issue models.Issue) error {
	return nil
}

func (R *Repo) UpdateIssue(issue models.Issue) error {
	return nil
}

func (R *Repo) CloseIssue(issueID int) error {
	return nil
}

func (R *Repo) GetIssuesList(repoID int) ([]models.Issue, error) {
	return []models.Issue{}, nil
}

func (R *Repo) CheckAccessIssue(userID, issueID int) (perm.Permission, error) {
	return perm.AdminAccess(), nil
}

func (R *Repo) CheckAccessRepo(userID, repoID int) (perm.Permission, error) {
	return perm.AdminAccess(), nil
}

func (R *Repo) GetIssues(issueID int) (models.Issue, error) {
	return models.Issue{}, nil
}
