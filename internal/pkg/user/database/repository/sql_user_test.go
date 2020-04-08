package repository

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Suite struct {
	suite.Suite
	repo UserRepo
	user models.User
	mock sqlmock.Sqlmock
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)
	s.user = models.User{
		Id:       1,
		Password: "123456789",
		Name:     "heehheheh",
		Login:    "kekmdda",
		Image:    "/image/test",
		Email:    "testik@email.test",
	}

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	require.NoError(s.T(), err)

	s.repo = NewUserRepo(sqlxDB, "/image/test", "", "")
}
func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}
func (s *Suite) TestGetUserByIdWithPass() {
	// good query
	rows := s.mock.
		NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(s.user.Id, s.user.Login, s.user.Email, s.user.Password, s.user.Name, s.user.Image)

	s.mock.
		ExpectQuery("SELECT").
		WithArgs(s.user.Login).
		WillReturnRows(rows)

	item, err := s.repo.GetUserByLoginWithPass(s.user.Login)
	if err != nil {
		s.T().Errorf("unexpected err: %s", err)
		return
	}
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	fmt.Println(item)
}
