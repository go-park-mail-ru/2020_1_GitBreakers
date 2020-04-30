package issues

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type issueTestSuite struct {
	suite.Suite
	issueMaker      models.User
	gitRepository   gitmodels.Repository
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

	s.gitRepository = gitmodels.Repository{
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

func (s *issueTestSuite) TestCreateIssuePositive() {
	for _, issue := range s.issues {
		s.mock.ExpectExec("INSERT INTO").WithArgs(
			issue.AuthorID,
			issue.RepoID,
			issue.Title,
			issue.Message,
			issue.Label,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.issueRepository.CreateIssue(issue)
		require.Nil(s.T(), err)
	}
}

func (s *issueTestSuite) TestCreateIssueNegative() {
	for _, issue := range s.issues {
		s.mock.ExpectExec("INSERT INTO").WithArgs(
			issue.AuthorID,
			issue.RepoID,
			issue.Title,
			issue.Message,
			issue.Label,
		).WillReturnError(sql.ErrConnDone)

		err := s.issueRepository.CreateIssue(issue)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.Is(err, sql.ErrConnDone))
	}
}

func (s *issueTestSuite) TestUpdatePositive() {
	for _, issue := range s.issues {
		s.mock.ExpectExec("UPDATE").WithArgs(
			issue.ID,
			issue.AuthorID,
			issue.RepoID,
			issue.Title,
			issue.Message,
			issue.Label,
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := s.issueRepository.UpdateIssue(issue)
		require.Nil(s.T(), err)
	}
}

func (s *issueTestSuite) TestUpdateNegative() {
	for _, issue := range s.issues {
		s.mock.ExpectExec("UPDATE").WithArgs(
			issue.ID,
			issue.AuthorID,
			issue.RepoID,
			issue.Title,
			issue.Message,
			issue.Label,
		).WillReturnError(sql.ErrConnDone)

		err := s.issueRepository.UpdateIssue(issue)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.Is(err, sql.ErrConnDone))
	}
}

func (s *issueTestSuite) TestCloseIssuePositive() {
	for _, issue := range s.issues {
		isRepoSuccessfullyClosed := s.mock.NewRows([]string{"result"})
		isRepoSuccessfullyClosed.AddRow(true)

		s.mock.ExpectQuery("UPDATE").
			WithArgs(issue.ID).
			WillReturnRows(isRepoSuccessfullyClosed)

		err := s.issueRepository.CloseIssue(issue.ID)
		require.Nil(s.T(), err)
	}
}

func (s *issueTestSuite) TestCloseIssueNegative() {
	for _, issue := range s.issues {
		s.mock.ExpectQuery("UPDATE").
			WithArgs(issue.ID).
			WillReturnError(sql.ErrNoRows)

		err := s.issueRepository.CloseIssue(issue.ID)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.Is(err, entityerrors.Invalid()))

		s.mock.ExpectQuery("UPDATE").
			WithArgs(issue.ID).
			WillReturnError(sql.ErrConnDone)

		err = s.issueRepository.CloseIssue(issue.ID)
		require.NotNil(s.T(), err)
		require.True(s.T(), errors.Is(err, sql.ErrConnDone))
	}
}

func (s *issueTestSuite) TestGetIssuesListPositive() {
	issuesRows := s.mock.NewRows(
		[]string{
			"id",
			"author_id",
			"repo_id",
			"title",
			"message",
			"label",
			"is_closed",
			"created_at",
		},
	)

	for _, issue := range s.issues {
		issuesRows.AddRow(
			issue.ID,
			issue.AuthorID,
			issue.RepoID,
			issue.Title,
			issue.Message,
			issue.Label,
			issue.IsClosed,
			issue.CreatedAt,
		)
	}

	offset := int64(0)
	limit := int64(len(s.issues))

	s.mock.ExpectQuery("SELECT").
		WithArgs(s.gitRepository.ID, limit, offset).
		WillReturnRows(issuesRows)

	result, err := s.issueRepository.GetOpenedIssuesList(s.gitRepository.ID, limit, offset)
	require.Nil(s.T(), err)
	require.EqualValues(s.T(), s.issues, result)
}

func (s *issueTestSuite) TestGetIssuesNegative() {
	offset := int64(0)
	limit := int64(len(s.issues))

	s.mock.ExpectQuery("SELECT").
		WithArgs(s.gitRepository.ID, limit, offset).
		WillReturnError(sql.ErrConnDone)

	_, err := s.issueRepository.GetOpenedIssuesList(s.gitRepository.ID, limit, offset)
	require.NotNil(s.T(), err)
}
