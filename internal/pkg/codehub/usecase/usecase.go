package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type UCCodeHub struct {
	Repo codehub.Repository
}

func (GD *UCCodeHub) ModifyStar(star models.Star) error {
	if star.Vote {
		return GD.Repo.AddStar(star.AuthorID, star.RepoID)
	} else {
		return GD.Repo.DelStar(star.AuthorID, star.RepoID)
	}
}

func (GD *UCCodeHub) GetStarredRepo(userID int) (models.RepoSet, error) {
	return GD.Repo.GetStarredRepo(userID)
}

func (GD *UCCodeHub) CreateIssue(issue models.Issue) error {
	permis, err := GD.Repo.CheckAccessRepo(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.Repo.CreateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) UpdateIssue(issue models.Issue) error {
	permis, err := GD.Repo.CheckAccessIssue(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.Repo.UpdateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}

}

func (GD *UCCodeHub) CloseIssue(issueID int, userID int) error {
	permis, err := GD.Repo.CheckAccessIssue(userID, issueID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.Repo.CloseIssue(issueID)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssuesList(repoID int, userID int) (models.IssuesSet, error) {
	permis, err := GD.Repo.CheckAccessRepo(userID, repoID)
	if err != nil {
		return nil, err
	}

	if permis != perm.NoAccess() {
		return GD.Repo.GetIssuesList(repoID)
	} else {
		return nil, entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssue(issueID int, userID int) (models.Issue, error) {
	permis, err := GD.Repo.CheckAccessIssue(userID, issueID)
	if err != nil {
		return models.Issue{}, err
	}

	if permis != perm.NoAccess() {
		return GD.Repo.GetIssues(issueID)
	} else {
		return models.Issue{}, entityerrors.AccessDenied()
	}
}
