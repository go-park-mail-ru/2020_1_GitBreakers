package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitModels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
	"time"
)

var gitTestRepoDir = os.TempDir() + "/" + uuid.NewV4().String()

type Suite struct {
	suite.Suite
	gitRepository    Repository
	gitRepoModel     gitModels.Repository
	gitRepoUserModel models.User
	mock             sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
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
	suite.Run(t, new(Suite))
	if err := os.RemoveAll(gitTestRepoDir); err != nil {
		fmt.Printf("error while removing test dir %s, err=%+v\n", gitTestRepoDir, err)
		t.Fail()
	}
}

func (s *Suite) TestRepositoryCreate() {
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
