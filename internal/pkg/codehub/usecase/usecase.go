package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type UCCodeHub struct {
	repo codehub.Repository
}

func (GD *UCCodeHub) ModifyStar(star models.Star) error {
	if star.Vote {
		return GD.repo.AddStar(star.AuthorID, star.RepoID)
	} else {
		return GD.repo.DelStar(star.AuthorID, star.RepoID)
	}
}

func (GD *UCCodeHub) GetStarredRepo(userID int) ([]gitmodels.Repository, error) {
	return GD.repo.GetStarredRepo(userID)
}

func (GD *UCCodeHub) CreateIssue(issue models.Issue) error {
	permis, err := GD.repo.CheckAccessRepo(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.repo.CreateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) UpdateIssue(issue models.Issue) error {
	permis, err := GD.repo.CheckAccessIssue(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.repo.UpdateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}

}

func (GD *UCCodeHub) CloseIssue(issueID int, userID int) error {
	permis, err := GD.repo.CheckAccessIssue(userID, issueID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.repo.CloseIssue(issueID)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssuesList(repoID int, userID int) ([]models.Issue, error) {
	permis, err := GD.repo.CheckAccessRepo(userID, repoID)
	if err != nil {
		return nil, err
	}

	if permis != perm.NoAccess() {
		return GD.repo.GetIssuesList(repoID)
	} else {
		return nil, entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssues(issueID int, userID int) (models.Issue, error) {
	permis, err := GD.repo.CheckAccessIssue(userID, issueID)
	if err != nil {
		return models.Issue{}, err
	}

	if permis != perm.NoAccess() {
		return GD.repo.GetIssues(issueID)
	} else {
		return models.Issue{}, entityerrors.AccessDenied()
	}
}
