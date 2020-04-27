package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
)

type UCCodeHub struct {
	RepoIssue codehub.RepoIssueI
	RepoStar  codehub.RepoStarI
	RepoNews  codehub.RepoNewsI
	GitRepo   git.Repository
	UserRepo  user.RepoUser
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
	repo, err := UC.GitRepo.GetByID(issue.RepoID)
	if err != nil {
		return entityerrors.DoesNotExist()
	}
	login, err := UC.UserRepo.GetLoginByID(repo.OwnerID)
	if err != nil {
		return entityerrors.DoesNotExist()
	}
	currUserId := issue.AuthorID

	permis, err := UC.GitRepo.CheckReadAccess(&currUserId, login, repo.Name)
	if err != nil {
		return err
	}

	if permis { //can create if repo not private
		return UC.RepoIssue.CreateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) UpdateIssue(issue models.Issue) error {
	permis, _ := GD.RepoIssue.CheckEditAccessIssue(issue.AuthorID, issue.RepoID)

	permis = perm.AdminAccess()
	if permis == perm.WriteAccess() || permis == perm.AdminAccess() || permis == perm.OwnerAccess() {
		return GD.RepoIssue.UpdateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}

}

func (GD *UCCodeHub) CloseIssue(issueID int64, userID int64) error {
	permis, _ := GD.RepoIssue.CheckEditAccessIssue(userID, issueID)
	permis = perm.AdminAccess()

	if permis == perm.WriteAccess() || permis == perm.AdminAccess() || permis == perm.OwnerAccess() {
		return GD.RepoIssue.CloseIssue(issueID)
	} else {
		return entityerrors.AccessDenied()
	}
}

func (UC *UCCodeHub) GetIssuesList(repoID, userID, limit, offset int64) (models.IssuesSet, error) {
	repo, err := UC.GitRepo.GetByID(repoID)
	if err != nil {
		return models.IssuesSet{}, entityerrors.DoesNotExist()
	}
	login, err := UC.UserRepo.GetLoginByID(repo.OwnerID)
	if err != nil {
		return models.IssuesSet{}, entityerrors.DoesNotExist()
	}

	permis, err := UC.GitRepo.CheckReadAccess(&userID, login, repo.Name)
	if err != nil {
		return models.IssuesSet{}, err
	}

	if permis {
		return UC.RepoIssue.GetIssuesList(repoID, limit, offset)
	} else {
		return nil, entityerrors.AccessDenied()
	}
}

func (GD *UCCodeHub) GetIssue(issueID int64, userID int64) (models.Issue, error) {
	return GD.RepoIssue.GetIssue(issueID)
}

func (UC *UCCodeHub) GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error) {
	repo, err := UC.GitRepo.GetByID(repoID)
	if err != nil {
		return models.NewsSet{}, entityerrors.DoesNotExist()
	}
	login, err := UC.UserRepo.GetLoginByID(repo.OwnerID)
	if err != nil {
		return models.NewsSet{}, entityerrors.DoesNotExist()
	}

	permis, err := UC.GitRepo.CheckReadAccess(&userID, login, repo.Name)
	if err != nil {
		return models.NewsSet{}, err
	}
	if permis {
		return UC.RepoNews.GetNews(repoID, limit, offset)
	} else {
		return models.NewsSet{}, entityerrors.AccessDenied()
	}
}
func (UC *UCCodeHub) GetUserStaredList(repoID int64, limit int64, offset int64) (models.UserSet, error) {
	return UC.RepoStar.GetUserStaredList(repoID, limit, offset)
}
