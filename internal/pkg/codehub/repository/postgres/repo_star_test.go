package postgres

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type starTestSuite struct {
	suite.Suite
	firstStargazer  models.User
	secondStargazer models.User
	repositories    models.RepoSet
	mock            sqlmock.Sqlmock
}

//IsExistStar(userID int64, repoID int64) (bool, error)
//AddStar(userID int64, repoID int64) error
//DelStar(userID int64, repoID int64) error
//GetStarredRepos(userID int64, limit int64, offset int64) ([]gitmodels.Repository, error)
//GetUserStaredList(repoID int64, limit int64, offset int64) ([]models.User, error)

func (s *starTestSuite) SetupSuite() {
	//var (
	//	db  *sql.DB
	//	err error
	//)

	s.firstStargazer = models.User{
		ID:        1,
		Password:  "",
		Name:      "first",
		Login:     "first",
		Image:     "",
		Email:     "first@mail.ru",
		CreatedAt: time.Time{},
	}

	s.secondStargazer = models.User{
		ID:        1,
		Password:  "",
		Name:      "second",
		Login:     "second",
		Image:     "",
		Email:     "second@mail.ru",
		CreatedAt: time.Time{},
	}

	s.repositories = models.RepoSet{
		gitmodels.Repository{
			ID:          1,
			OwnerID:     1,
			Name:        "first",
			Description: "first",
			IsFork:      false,
			CreatedAt:   time.Time{},
			IsPublic:    false,
			Stars:       0,
		},
		gitmodels.Repository{
			ID:          2,
			OwnerID:     2,
			Name:        "second",
			Description: "second",
			IsFork:      false,
			CreatedAt:   time.Time{},
			IsPublic:    false,
			Stars:       0,
		},
	}

}

func TestInit(t *testing.T) {
	suite.Run(t, new(starTestSuite))
}

func (s *starTestSuite) TestIsExistStar() {

}

func (s *starTestSuite) TestDelStar() {

}

func (s *starTestSuite) TestGetStarredRepos() {

}

func (s *starTestSuite) TestGetUserStarredList() {

}
