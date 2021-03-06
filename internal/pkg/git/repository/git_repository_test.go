package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitModels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
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
		ID:       1,
		Password: "123456789",
		Name:     "heehheheh",
		Login:    "kekmdda",
		Image:    "/image/test",
		Email:    "testik@email.test",
	}
	s.gitRepoModel = gitModels.Repository{
		ID:                   1,
		OwnerID:              s.gitRepoUserModel.ID,
		Name:                 "test_repo",
		Description:          "test repository",
		IsFork:               false,
		CreatedAt:            time.Now(),
		IsPublic:             false,
		ParentRepositoryInfo: gitModels.ParentRepositoryInfo{},
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
	repoRow.AddRow(repo.ID)

	userRow := s.mock.NewRows([]string{"login"})
	userRow.AddRow(user.Login)

	s.mock.ExpectBegin()

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(repo.OwnerID, repo.Name).WillReturnRows(repoExitsInDbRow)

	s.mock.ExpectQuery("INSERT").
		WithArgs(
			repo.OwnerID,
			repo.Name,
			repo.Description,
			repo.IsPublic,
			repo.IsFork,
			repo.ParentRepositoryInfo.ID,
		).WillReturnRows(repoRow)

	s.mock.ExpectExec("INSERT").WithArgs(repo.OwnerID, repo.ID, perm.OwnerAccess()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	s.mock.ExpectQuery("SELECT").WithArgs(user.ID).WillReturnRows(userRow)
	s.mock.ExpectCommit()

	id, err := s.gitRepository.Create(repo)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), id, repo.ID)
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
		WithArgs(repo.OwnerID, repo.Name).WillReturnRows(repoExitsInDbRow)

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
		WithArgs(user.Login, repo.Name).
		WillReturnRows(haveReadAccessDbRow)

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user.Login, repo.Name, user.ID).
		WillReturnRows(haveReadAccessDbRow)

	haveReadAccess, err := s.gitRepository.CheckReadAccess(&user.ID, user.Login, repo.Name)
	require.Error(s.T(), err)
	require.EqualValues(s.T(), haveReadAccess, false)
}

func (s *gitRepoTestSuite) TestCheckReadAccessNegative() {
	repo := s.gitRepoModel
	user := s.gitRepoUserModel

	haveReadAccessDbRow := s.mock.NewRows([]string{"exists"})
	haveReadAccessDbRow.AddRow(false)

	s.mock.ExpectQuery("SELECT EXISTS").
		WithArgs(user.Login, repo.Name).WillReturnRows(haveReadAccessDbRow)

	haveReadAccess, err := s.gitRepository.CheckReadAccess(&user.ID, user.Login, repo.Name)
	require.EqualValues(s.T(), entityerrors.DoesNotExist(), err)
	require.EqualValues(s.T(), haveReadAccess, false)
}

func (s *gitRepoTestSuite) TestGetById() {
	repo := s.gitRepoModel

	repoDbRow := s.mock.NewRows([]string{
		"id",
		"owner_id",
		"name",
		"description",
		"is_fork",
		"is_public",
		"stars",
		"forks",
		"merge_requests_open",
		"created_at",
		"user_login",
		"parent_id",
		"parent_owner_id",
		"parent_name",
		"parent_user_login",
	})
	repoDbRow.AddRow(
		repo.ID,
		repo.OwnerID,
		repo.Name,
		repo.Description,
		repo.IsFork,
		repo.IsPublic,
		repo.Stars,
		repo.Forks,
		repo.MergeRequestsOpen,
		repo.CreatedAt,
		repo.AuthorLogin,
		repo.ParentRepositoryInfo.ID,
		repo.ParentRepositoryInfo.OwnerID,
		repo.ParentRepositoryInfo.Name,
		repo.ParentRepositoryInfo.AuthorLogin,
	)

	s.mock.ExpectQuery("SELECT").
		WithArgs(repo.ID).WillReturnRows(repoDbRow)

	repoModel, err := s.gitRepository.GetByID(repo.ID)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), repoModel, repo)
}

func (s *gitRepoTestSuite) TestGetByIdNegativeDoesNotExist() {
	repo := s.gitRepoModel

	s.mock.ExpectQuery("SELECT").
		WithArgs(repo.ID).WillReturnError(sql.ErrNoRows)

	_, err := s.gitRepository.GetByID(repo.ID)

	require.NotNil(s.T(), err)
	require.EqualValues(s.T(), err, entityerrors.DoesNotExist())
}
