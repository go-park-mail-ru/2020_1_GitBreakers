package postgres

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/jmoiron/sqlx"
)

type RepoStar struct {
	DB *sqlx.DB
}

func NewRepoStar(db *sqlx.DB) codehub.RepoStarI {
	return &RepoStar{DB: db}
}
func (R *RepoStar) AddStar(userID int64, repoID int64) error {
	panic("implement me")
}

func (R *RepoStar) DelStar(userID int64, repoID int64) error {
	panic("implement me")
}

func (R *RepoStar) GetStarredRepo(userID int64, limit int64, offset int64) ([]gitmodels.Repository, error) {
	panic("implement me")
}
