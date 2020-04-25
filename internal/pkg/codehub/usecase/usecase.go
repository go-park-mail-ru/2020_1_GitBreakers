package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type UCCodeHub struct {
	RepoIssue codehub.RepoIssueI
	RepoStar  codehub.RepoStarI
	RepoNews  codehub.RepoNewsI
}

func (GD *UCCodeHub) ModifyStar(star models.Star) error {
	if star.Vote {
		return GD.RepoStar.AddStar(star.AuthorID, star.RepoID)
	} else {
		return GD.RepoStar.DelStar(star.AuthorID, star.RepoID)
	}
}

func (GD *UCCodeHub) GetStarredRepo(userID int64) (models.RepoSet, error) {
	return GD.RepoStar.GetStarredRepo(userID)
}

func (GD *UCCodeHub) CreateIssue(issue models.Issue) error {
	permis, err := GD.RepoIssue.CheckAccessRepo(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis != perm.NoAccess() { //can create if repo not private
		return GD.RepoIssue.CreateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) UpdateIssue(issue models.Issue) error {
	permis, err := GD.RepoIssue.CheckAccessIssue(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.RepoIssue.UpdateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}

}

func (GD *UCCodeHub) CloseIssue(issueID int64, userID int64) error {
	permis, err := GD.RepoIssue.CheckAccessIssue(userID, issueID)
	if err != nil {
		return err
	}

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() {
		return GD.RepoIssue.CloseIssue(issueID)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssuesList(repoID int64, userID int64) (models.IssuesSet, error) {
	permis, err := GD.RepoIssue.CheckAccessRepo(userID, repoID)
	if err != nil {
		return nil, err
	}

	if permis != perm.NoAccess() {
		return GD.RepoIssue.GetIssuesList(repoID)
	} else {
		return nil, entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssue(issueID int64, userID int64) (models.Issue, error) {
	permis, err := GD.RepoIssue.CheckAccessIssue(userID, issueID)
	if err != nil {
		return models.Issue{}, err
	}

	if permis != perm.NoAccess() {
		return GD.RepoIssue.GetIssues(issueID)
	} else {
		return models.Issue{}, entityerrors.AccessDenied()
	}
}
func (GD *UCCodeHub) GetNews(repoID int64, userID int64) (models.Issue, error) {
	//permis, err := GD.RepoIssue.CheckAccessRepo(userID, repoID)

}
