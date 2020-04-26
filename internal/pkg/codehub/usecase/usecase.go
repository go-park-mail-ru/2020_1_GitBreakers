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

func (UC *UCCodeHub) ModifyStar(star models.Star) error {
	if star.Vote {
		return UC.RepoStar.AddStar(star.AuthorID, star.RepoID)
	} else {
		return UC.RepoStar.DelStar(star.AuthorID, star.RepoID)
	}
}

func (UC *UCCodeHub) GetStarredRepos(userID, limit, offset int64) (models.RepoSet, error) {
	return UC.RepoStar.GetStarredRepos(userID, limit, offset)
}

func (UC *UCCodeHub) CreateIssue(issue models.Issue) error {
	permis, err := UC.RepoIssue.CheckAccessRepo(issue.AuthorID, issue.RepoID)
	if err != nil {
		return err
	}

	if permis != perm.NoAccess() { //can create if repo not private
		return UC.RepoIssue.CreateIssue(issue)
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

func (GD *UCCodeHub) GetIssuesList(repoID, userID, limit, offset int64) (models.IssuesSet, error) {
	permis, err := GD.RepoIssue.CheckAccessRepo(userID, repoID)
	if err != nil {
		return nil, err
	}

	if permis != perm.NoAccess() {
		return GD.RepoIssue.GetIssuesList(repoID, limit, offset)
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
		return GD.RepoIssue.GetIssue(issueID)
	} else {
		return models.Issue{}, entityerrors.AccessDenied()
	}
}

func (UC *UCCodeHub) GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error) {
	permis, err := UC.RepoIssue.CheckAccessRepo(userID, repoID)
	if err != nil {
		return nil, err
	}

	if permis != perm.NoAccess() {
		return UC.RepoNews.GetNews(repoID, limit, offset)
	} else {
		return models.NewsSet{}, entityerrors.AccessDenied()
	}
}
func (UC *UCCodeHub) GetUserStaredList(repoID int64, limit int64, offset int64) (models.UserSet, error) {
	return UC.RepoStar.GetUserStaredList(repoID, limit, offset)
}
