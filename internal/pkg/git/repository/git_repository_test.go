package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitModels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

var gitTestRepoDir = os.TempDir() + "/" + uuid.NewV4().String()

type gitRepoTestSuite struct {
	suite.Suite
	gitRepository    Repository
	gitRepoModel     gitModels.Repository
	gitRepoUserModel models.User
	mock             sqlmock.Sqlmock
}

func (s *gitRepoTestSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	s.gitRepoUserModel = models.User{
		Id:       1,
		Password: "123456789",
		Name:     "heehheheh",
		Login:    "kekmdda",
		Image:    "/image/test",
		Email:    "testik@email.test",
	}
	s.gitRepoModel = gitModels.Repository{
		Id:          1,
		OwnerId:     s.gitRepoUserModel.Id,
		Name:        "test_repo",
		Description: "test repository",
		IsFork:      false,
		CreatedAt:   time.Now(),
		IsPublic:    false,
	}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	require.NoError(s.T(), err)

	s.gitRepository = NewRepository(sqlxDB, gitTestRepoDir)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(gitRepoTestSuite))
	if err := os.RemoveAll(gitTestRepoDir); err != nil {
		fmt.Printf("error while removing test dir %s, err=%+v\n", gitTestRepoDir, err)
		t.Fail()
	}
}

func (s *gitRepoTestSuite) TestRepositoryCreate() {
	repo := s.gitRepoModel
	user := s.gitRepoUserModel

	defer func() {
		if err := os.RemoveAll(gitTestRepoDir + "/" + user.Name); err != nil {
			s.Failf("error while removing repo test dir %s, err=%+v\n", user.Name, err)
		}
	}()

	repoExitsInDbRow := s.mock.NewRows([]string{"exists"})
	repoExitsInDbRow.AddRow(false)

	repoRow := s.mock.NewRows([]string{"id"})
	repoRow.AddRow(repo.Id)

	userRow := s.mock.NewRows([]string{"login"})
	userRow.AddRow(user.Login)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(repo.OwnerId, repo.Name).WillReturnRows(repoExitsInDbRow)

	s.mock.ExpectQuery("INSERT").
		WithArgs(
			repo.OwnerId,
			repo.Name,
			repo.Description,
			repo.IsPublic,
			repo.IsFork,
		).WillReturnRows(repoRow)

	s.mock.ExpectExec("INSERT").WithArgs(repo.OwnerId, repo.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("SELECT").WithArgs(user.Id).WillReturnRows(userRow)
	s.mock.ExpectCommit()

	id, err := s.gitRepository.Create(repo)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), id, repo.Id)
}

func (s *gitRepoTestSuite) TestRepositoryCreateNegativeRepoExistInDb() {
	repo := s.gitRepoModel
	user := s.gitRepoUserModel

	defer func() {
		if err := os.RemoveAll(gitTestRepoDir + "/" + user.Name); err != nil {
			s.Failf("error while removing repo test dir %s, err=%+v\n", user.Name, err)
		}
	}()

	repoExitsInDbRow := s.mock.NewRows([]string{"exists"})
	repoExitsInDbRow.AddRow(true)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(repo.OwnerId, repo.Name).WillReturnRows(repoExitsInDbRow)

	s.mock.ExpectRollback()

	_, err := s.gitRepository.Create(repo)
	require.NotNil(s.T(), errors.Cause(err), entityerrors.AlreadyExist())
}

func (s *gitRepoTestSuite) TestCheckReadAccess() {
	repo := s.gitRepoModel
	user := s.gitRepoUserModel

	haveReadAccessDbRow := s.mock.NewRows([]string{"exists"})
	haveReadAccessDbRow.AddRow(true)

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user.Login, repo.Name, user.Id).WillReturnRows(haveReadAccessDbRow)

	haveReadAccess, err := s.gitRepository.CheckReadAccess(&user.Id, user.Login, repo.Name)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), haveReadAccess, true)
}

func (s *gitRepoTestSuite) TestCheckReadAccessNegative() {
	repo := s.gitRepoModel
	user := s.gitRepoUserModel

	haveReadAccessDbRow := s.mock.NewRows([]string{"exists"})
	haveReadAccessDbRow.AddRow(false)

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user.Login, repo.Name, user.Id).WillReturnRows(haveReadAccessDbRow)

	haveReadAccess, err := s.gitRepository.CheckReadAccess(&user.Id, user.Login, repo.Name)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), haveReadAccess, false)
}

func (s *gitRepoTestSuite) TestGetById() {
	repo := s.gitRepoModel

	repoDbRow := s.mock.NewRows([]string{
		"id",
		"owner_id",
		"name",
		"description",
		"is_public",
		"is_fork",
		"created_at",
	})
	repoDbRow.AddRow(
		repo.Id,
		repo.OwnerId,
		repo.Name,
		repo.Description,
		repo.IsPublic,
		repo.IsFork,
		repo.CreatedAt,
	)

	s.mock.ExpectQuery("SELECT").
		WithArgs(repo.Id).WillReturnRows(repoDbRow)

	repoModel, err := s.gitRepository.GetById(repo.Id)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), repoModel, repo)
}

func (s *gitRepoTestSuite) TestGetByIdNegativeDoesNotExist() {
	repo := s.gitRepoModel

	s.mock.ExpectQuery("SELECT").
		WithArgs(repo.Id).WillReturnError(sql.ErrNoRows)

	_, err := s.gitRepository.GetById(repo.Id)

	require.NotNil(s.T(), err)
	require.EqualValues(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}
