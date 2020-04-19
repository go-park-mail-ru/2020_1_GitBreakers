package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
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
	return GD.repo.CreateIssue(issue)
}

func (GD *UCCodeHub) UpdateIssue(issue models.Issue) error {
	return GD.repo.UpdateIssue(issue)
}

func (GD *UCCodeHub) CloseIssue(issueID int) error {
	return GD.repo.CloseIssue(issueID)
}

func (GD *UCCodeHub) GetIssuesList(repoID int) ([]models.Issue, error) {
	return GD.repo.GetIssuesList(repoID)
}
