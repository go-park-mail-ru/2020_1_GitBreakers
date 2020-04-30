package issues

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type issueTestSuite struct {
	suite.Suite
	issueMaker      models.User
	repository      gitmodels.Repository
	issues          models.IssuesSet
	issueRepository IssueRepository
	mock            sqlmock.Sqlmock
}

func (s *issueTestSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	s.issueMaker = models.User{
		ID:        1,
		Password:  "",
		Name:      "issue maker",
		Login:     "issue_maker",
		Image:     "",
		Email:     "issue_maker@mail.ru",
		CreatedAt: time.Time{},
	}

	s.repository = gitmodels.Repository{
		ID:          1,
		OwnerID:     1,
		Name:        "testrepo",
		Description: "testrepo descr",
		IsFork:      false,
		CreatedAt:   time.Time{},
		IsPublic:    false,
		Stars:       0,
	}

	s.issues = models.IssuesSet{
		models.Issue{
			ID:        1,
			AuthorID:  1,
			RepoID:    1,
			Title:     "first issue",
			Message:   "first issue",
			Label:     "bug",
			IsClosed:  false,
			CreatedAt: time.Time{},
		},
		models.Issue{
			ID:        2,
			AuthorID:  1,
			RepoID:    1,
			Title:     "second issue",
			Message:   "second issue",
			Label:     "feature",
			IsClosed:  false,
			CreatedAt: time.Time{},
		},
	}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	require.NoError(s.T(), err)

	s.issueRepository = NewIssueRepository(sqlxDB)
}

func TestInit(t *testing.T) {
	suite.Run(t, new(issueTestSuite))
}
