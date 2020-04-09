package repository

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"reflect"
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
func (s *Suite) TestGetUserByLoginWithPass() {
	user := s.user
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(s.T(), err)
	user.Password = string(hash)
	rows := s.mock.
		NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	UserFromDB, err := s.repo.GetUserByLoginWithPass(user.Login)
	require.Nil(s.T(), err)

	if !reflect.DeepEqual(UserFromDB, user) {
		s.Assert()
	}

	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.GetUserByLoginWithPass(user.Login)
	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())

}
func (s *Suite) TestGetUserByLoginWithoutPass() {
	user := s.user

	rows := s.mock.
		NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.
		ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	UserFromDB, err := s.repo.GetByLoginWithoutPass(user.Login)

	require.NotEqual(s.T(), UserFromDB.Password, user.Password)
	require.Equal(s.T(), UserFromDB.Password, "")
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	someErr := errors.New("some db error")
	s.mock.ExpectQuery("SELECT").WithArgs(user.Login).
		WillReturnError(someErr)

	_, err = s.repo.GetUserByLoginWithPass(user.Login)
	if reflect.DeepEqual(UserFromDB, user) {
		s.Assert()
	}

	require.Equal(s.T(), errors.Cause(err), someErr)
}
func (s *Suite) TestGetUserByIdWithoutPass() {
	user := s.user

	rows := s.mock.
		NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.
		ExpectQuery("SELECT").
		WithArgs(user.Id).
		WillReturnRows(rows)

	UserFromDB, err := s.repo.GetUserByIdWithoutPass(user.Id)

	require.NotEqual(s.T(), UserFromDB.Password, user.Password)
	require.Equal(s.T(), UserFromDB.Password, "")
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	someErr := errors.New("some db error")
	s.mock.ExpectQuery("SELECT").WithArgs(user.Id).
		WillReturnError(someErr)

	_, err = s.repo.GetUserByIdWithoutPass(user.Id)
	if reflect.DeepEqual(UserFromDB, user) {
		s.Assert()
	}

	require.Equal(s.T(), errors.Cause(err), someErr)
}

func (s *Suite) TestGetUserByIdWithPass() {
	user := s.user
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	require.NoError(s.T(), err)
	user.Password = string(hash)
	rows := s.mock.
		NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.
		ExpectQuery("SELECT").
		WithArgs(user.Id).
		WillReturnRows(rows)

	UserFromDB, err := s.repo.GetUserByIdWithPass(user.Id)
	require.Nil(s.T(), err)

	if !reflect.DeepEqual(UserFromDB, user) {
		s.Assert()
	}

	s.mock.ExpectQuery("SELECT").WithArgs(user.Id).
		WillReturnError(sql.ErrNoRows)
	_, err = s.repo.GetUserByIdWithPass(user.Id)
	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())

}
