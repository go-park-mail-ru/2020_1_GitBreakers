package postgres

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/jmoiron/sqlx"
)

type RepoIssue struct {
	DB *sqlx.DB
}

func NewRepoIssue(db *sqlx.DB) RepoIssue {
	return RepoIssue{
		DB: db,
	}
}

func (R *RepoIssue) CreateIssue(issue models.Issue) error {
	panic("implement me")
}

func (R *RepoIssue) UpdateIssue(issue models.Issue) error {
	panic("implement me")
}

func (R *RepoIssue) CloseIssue(issueID int64) error {
	panic("implement me")
}

func (R *RepoIssue) GetIssuesList(repoID int64, limit int64, offset int64) ([]models.Issue, error) {
	panic("implement me")
}

func (R *RepoIssue) CheckAccessIssue(userID, issueID int64) (perm.Permission, error) {
	panic("implement me")
}

func (R *RepoIssue) CheckAccessRepo(userID, repoID int64) (perm.Permission, error) {
	panic("implement me")
}

func (R *RepoIssue) GetIssue(issueID int64) (models.Issue, error) {
	panic("implement me")
}
