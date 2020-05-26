package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/pkg/errors"
)

type UCCodeHub struct {
	RepoIssue  codehub.RepoIssueI
	RepoStar   codehub.RepoStarI
	RepoNews   codehub.RepoNewsI
	GitRepo    git.GitRepoI
	UserRepo   user.RepoUser
	SearchRepo codehub.RepoSearchI
	RepoMerge  codehub.RepoMergeI
}

func (UC *UCCodeHub) ModifyStar(star models.Star) error {
	if star.Vote {
		return UC.RepoStar.AddStar(star.AuthorID, star.RepoID)
	} else {
		return UC.RepoStar.DelStar(star.AuthorID, star.RepoID)
	}
}

func (UC *UCCodeHub) GetStarredRepos(userID, limit, offset int64) (models.RepoSet, error) {
	repolist, err := UC.RepoStar.GetStarredRepos(userID, limit, offset)
	if repolist != nil {
		return repolist, err
	}
	return models.RepoSet{}, err
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
	permis, _ := GD.RepoIssue.CheckEditAccessIssue(issue.AuthorID, issue.ID)
	//todo check that work
	if permis == perm.WriteAccess() || permis == perm.AdminAccess() || permis == perm.OwnerAccess() {
		return GD.RepoIssue.UpdateIssue(issue)
	} else {
		return entityerrors.AccessDenied()
	}

}

func (GD *UCCodeHub) CloseIssue(issueID int64, userID int64) error {
	permis, _ := GD.RepoIssue.CheckEditAccessIssue(userID, issueID)

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
		return models.IssuesSet{}, entityerrors.AccessDenied()
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
	UserSet, err := UC.RepoStar.GetUserStaredList(repoID, limit, offset)
	if UserSet != nil {
		return UserSet, err
	}
	return models.UserSet{}, err
}
func (UC *UCCodeHub) Search(query, params string, limit, offset, userID int64) (interface{}, error) {
	switch params {
	case "allusers":
		return UC.SearchRepo.GetFromUsers(query, limit, offset)

	case "allrepo":
		return UC.SearchRepo.GetFromAllRepos(query, limit, offset)

	case "myrepo":
		return UC.SearchRepo.GetFromOwnRepos(query, limit, offset, userID)

	case "starredrepo":
		return UC.SearchRepo.GetFromStarredRepos(query, limit, offset, userID)

	default:
		return nil, entityerrors.Invalid()
	}
}
func (UC *UCCodeHub) CreatePL(request models.PullRequest) (models.PullRequest, error) {
	if request.FromRepoID == nil {
		return models.PullRequest{}, entityerrors.Invalid()
	}

	toRepoIdAccess, err := UC.GitRepo.CheckReadAccessById(request.AuthorId, request.ToRepoID)
	if err != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}
	if !toRepoIdAccess {
		return models.PullRequest{}, entityerrors.AccessDenied()
	}

	permission, err := UC.GitRepo.GetPermissionByID(request.AuthorId, *request.FromRepoID)
	if err != nil {
		return models.PullRequest{}, errors.WithStack(err)
	}

	if permission == perm.OwnerAccess() || permission == perm.AdminAccess() {
		return UC.RepoMerge.CreateMR(request)
	}

	return models.PullRequest{}, entityerrors.AccessDenied()
}

func (UC *UCCodeHub) GetPLIn(repo gitmodels.Repository, limit int64, offset int64) (models.PullReqSet, error) {
	_, err := UC.GitRepo.GetByID(repo.ID)
	if err != nil {
		return models.PullReqSet{}, err
	}
	return UC.RepoMerge.GetAllMRIn(repo.ID, limit, offset)
}

func (UC *UCCodeHub) GetPLOut(repo gitmodels.Repository, limit int64, offset int64) (models.PullReqSet, error) {
	return UC.RepoMerge.GetAllMROut(repo.ID, limit, offset)
}

func (UC *UCCodeHub) ApprovePL(pl models.PullRequest, userID int64) error {
	permission, err := UC.GitRepo.GetPermissionByID(&userID, pl.ToRepoID)
	if err != nil {
		return errors.WithStack(err)
	}

	if permission == perm.OwnerAccess() || permission == perm.AdminAccess() {
		approver, err := UC.UserRepo.GetUserByIDWithoutPass(userID)
		if err != nil {
			return errors.WithStack(err)
		}

		return UC.RepoMerge.ApproveMerge(pl.ID, approver)
	}

	return entityerrors.AccessDenied()
}

func (UC *UCCodeHub) ClosePL(pl models.PullRequest, userID int64) error {
	var isCorrect bool
	if pl.AuthorId != nil && *pl.AuthorId == userID {
		isCorrect = true
	} else {
		permission, err := UC.GitRepo.GetPermissionByID(&userID, pl.ToRepoID)
		if err != nil {
			return errors.WithStack(err)
		}

		if permission == perm.OwnerAccess() || permission == perm.AdminAccess() {
			isCorrect = true
		}
	}

	if isCorrect {
		return UC.RepoMerge.RejectMR(userID, pl.ID)
	}

	return entityerrors.AccessDenied()
}

func (UC *UCCodeHub) GetAllMRUser(userID, limit, offset int64) (models.PullReqSet, error) {
	MRList, err := UC.RepoMerge.GetAllMRForUser(userID, limit, offset)
	return MRList, errors.WithStack(err)
}

func (UC *UCCodeHub) GetMRByID(mrID int64) (models.PullRequest, error) {
	return UC.RepoMerge.GetMRByID(mrID)
}

func (UC *UCCodeHub) GetMRDiffByID(mrID int64) (string, error) {
	return UC.RepoMerge.GetMRDiffByID(mrID)
}
