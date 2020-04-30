package stars

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type starTestSuite struct {
	suite.Suite
	firstStargazer  models.User
	secondStargazer models.User
	gitRepositories models.RepoSet
	starRepository  StarRepository
	mock            sqlmock.Sqlmock
}

func (s *starTestSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

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

	s.gitRepositories = models.RepoSet{
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

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	require.NoError(s.T(), err)

	s.starRepository = NewStarRepository(sqlxDB)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(starTestSuite))
}

func (s *starTestSuite) testIsExistStarPositive(expectedResult bool) {
	repoExitsInDbRow := s.mock.NewRows([]string{"exists"})
	repoExitsInDbRow.AddRow(expectedResult)

	firstRepo := s.gitRepositories[0]

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(firstRepo.ID, firstRepo.OwnerID).WillReturnRows(repoExitsInDbRow)

	result, err := s.starRepository.IsExistStar(firstRepo.ID, firstRepo.OwnerID)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), expectedResult, result)
}

func (s *starTestSuite) TestIsExistStarPositive() {
	s.testIsExistStarPositive(true)
}

func (s *starTestSuite) TestIsExistStarNegative() {
	s.testIsExistStarPositive(false)
}

func (s *starTestSuite) TestAddStarPositive() {
	repoExitsInDbRow := s.mock.NewRows([]string{"exists"})
	repoExitsInDbRow.AddRow(false)

	firstRepo := s.gitRepositories[0]

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(firstRepo.ID, firstRepo.OwnerID).WillReturnRows(repoExitsInDbRow)

	s.mock.ExpectExec("INSERT INTO").
		WithArgs(s.firstStargazer.ID, firstRepo.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.starRepository.AddStar(s.firstStargazer.ID, firstRepo.ID)
	require.Nil(s.T(), err)
}

func (s *starTestSuite) TestAddStarNegative() {
	repoExitsInDbRow := s.mock.NewRows([]string{"exists"})
	repoExitsInDbRow.AddRow(true)

	firstRepo := s.gitRepositories[0]

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(firstRepo.ID, firstRepo.OwnerID).WillReturnRows(repoExitsInDbRow)

	err := s.starRepository.AddStar(s.firstStargazer.ID, firstRepo.ID)
	require.NotNil(s.T(), err)
	require.EqualValues(s.T(), entityerrors.AlreadyExist(), err)
}

func (s *starTestSuite) TestDelStarPositive() {
	isRepoSuccessfullyDeleted := s.mock.NewRows([]string{"exists"})
	isRepoSuccessfullyDeleted.AddRow(true)

	firstRepo := s.gitRepositories[0]

	s.mock.ExpectQuery("DELETE FROM").
		WithArgs(s.firstStargazer.ID, firstRepo.ID).
		WillReturnRows(isRepoSuccessfullyDeleted)

	err := s.starRepository.DelStar(s.firstStargazer.ID, firstRepo.ID)
	require.Nil(s.T(), err)
}

func (s *starTestSuite) TestDelStarNegative() {
	firstRepo := s.gitRepositories[0]

	s.mock.ExpectQuery("DELETE FROM").
		WithArgs(s.firstStargazer.ID, firstRepo.ID).
		WillReturnError(sql.ErrNoRows)

	err := s.starRepository.DelStar(s.firstStargazer.ID, firstRepo.ID)
	require.NotNil(s.T(), err)
	require.EqualValues(s.T(), entityerrors.DoesNotExist(), err)
}

func (s *starTestSuite) TestGetStarredRepos() {
	starredRepositories := s.mock.NewRows(
		[]string{
			"id",
			"owner_id",
			"name",
			"description",
			"is_fork",
			"created_at",
			"is_public",
			"stars",
		},
	)

	firstRepo := s.gitRepositories[0]
	secondRepo := s.gitRepositories[1]

	starredRepositories.AddRow(
		firstRepo.ID,
		firstRepo.OwnerID,
		firstRepo.Name,
		firstRepo.Description,
		firstRepo.IsFork,
		firstRepo.CreatedAt,
		firstRepo.IsPublic,
		firstRepo.Stars,
	)

	starredRepositories.AddRow(
		secondRepo.ID,
		secondRepo.OwnerID,
		secondRepo.Name,
		secondRepo.Description,
		secondRepo.IsFork,
		secondRepo.CreatedAt,
		secondRepo.IsPublic,
		secondRepo.Stars,
	)

	var limit int64 = 3
	var offset int64 = 0

	s.mock.ExpectQuery("SELECT").
		WithArgs(s.firstStargazer.ID, limit, offset).
		WillReturnRows(starredRepositories)

	result, err := s.starRepository.GetStarredRepos(s.firstStargazer.ID, limit, offset)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), s.gitRepositories, result)
}

func (s *starTestSuite) TestGetUserStarredList() {
	starredRepositories := s.mock.NewRows(
		[]string{
			"id",
			"login",
			"email",
			"name",
			"avatar_path",
			"created_at",
		},
	)

	starredRepositories.AddRow(
		s.firstStargazer.ID,
		s.firstStargazer.Login,
		s.firstStargazer.Email,
		s.firstStargazer.Name,
		s.firstStargazer.Image,
		s.firstStargazer.CreatedAt,
	)

	starredRepositories.AddRow(
		s.secondStargazer.ID,
		s.secondStargazer.Login,
		s.secondStargazer.Email,
		s.secondStargazer.Name,
		s.secondStargazer.Image,
		s.secondStargazer.CreatedAt,
	)

	var limit int64 = 3
	var offset int64 = 0

	firstRepository := s.gitRepositories[0]

	s.mock.ExpectQuery("SELECT").
		WithArgs(firstRepository.ID, limit, offset).
		WillReturnRows(starredRepositories)

	result, err := s.starRepository.GetUserStaredList(firstRepository.ID, limit, offset)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), []models.User{s.firstStargazer, s.secondStargazer}, result)
}
