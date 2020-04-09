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

func (s *Suite) SetupTest() {
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

	someErr := errors.New("some db error")
	s.mock.ExpectQuery("SELECT").WithArgs(user.Id).
		WillReturnError(someErr)

	_, err = s.repo.GetUserByIdWithoutPass(user.Id)
	if reflect.DeepEqual(UserFromDB, user) {
		s.Assert()
	}

	require.Equal(s.T(), errors.Cause(err), someErr)
}

//
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

//
func (s *Suite) TestIsExists() {
	user := s.user
	rows := s.mock.NewRows([]string{"is_exsist"})
	rows.AddRow(true)
	s.mock.ExpectQuery("SELECT  ").
		WithArgs(user.Login, user.Email).
		WillReturnRows(rows)
	isExsist, err := s.repo.IsExists(user)
	require.Nil(s.T(), err)
	require.Equal(s.T(), isExsist, true)

}

//
func (s *Suite) TestUpdate() {
	user := s.user
	rows := s.mock.NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.ExpectExec("UPDATE").
		WithArgs(user.Id, user.Email, user.Name, user.Image, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Update(user)
	require.Nil(s.T(), err)
	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	s.mock.ExpectExec("UPDATE").
		WithArgs(user.Id, user.Email, user.Name, user.Image, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 0))

	err = s.repo.Update(user)

	require.Equal(s.T(), err, entityerrors.Invalid())

}

//
func (s *Suite) TestUserCanUpdate() {
	user := s.user
	rows := s.mock.NewRows([]string{""})
	rows.AddRow(1)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Login, user.Email).
		WillReturnRows(rows)

	isCanUpdate, err := s.repo.UserCanUpdate(user)

	require.Nil(s.T(), err)
	require.True(s.T(), isCanUpdate)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

	rows = s.mock.NewRows([]string{""})
	rows.AddRow(2)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Login, user.Email).
		WillReturnRows(rows)

	isCanUpdate, err = s.repo.UserCanUpdate(user)

	require.Nil(s.T(), err, )
	require.False(s.T(), isCanUpdate)

	if err := s.mock.ExpectationsWereMet(); err != nil {
		s.T().Errorf("there were unfulfilled expectations: %s", err)
	}

}

//
func (s *Suite) TestGetLoginById() {
	user := s.user
	rows := s.mock.NewRows([]string{"login"})
	rows.AddRow(user.Login)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Id).
		WillReturnRows(rows)

	loginFromDB, err := s.repo.GetLoginByID(user.Id)

	require.Nil(s.T(), err)
	require.Equal(s.T(), loginFromDB, user.Login)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Id).
		WillReturnError(sql.ErrNoRows)

	loginFromDB, err = s.repo.GetLoginByID(user.Id)

	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Id).
		WillReturnError(errors.New("some errors"))

	loginFromDB, err = s.repo.GetLoginByID(user.Id)

	require.NotEqual(s.T(), errors.Cause(err), entityerrors.DoesNotExist())

}

//
func (s *Suite) TestGetIdByLogin() {
	user := s.user
	rows := s.mock.NewRows([]string{"id"})
	rows.AddRow(user.Id)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	loginFromDB, err := s.repo.GetIdByLogin(user.Login)

	require.Nil(s.T(), err)
	require.Equal(s.T(), loginFromDB, user.Id)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnError(sql.ErrNoRows)

	loginFromDB, err = s.repo.GetIdByLogin(user.Login)

	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Id).
		WillReturnError(errors.New("some errors"))

	loginFromDB, err = s.repo.GetIdByLogin(user.Login)

	require.NotEqual(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}

func (s *Suite) TestCreate() {
	user := s.user
	rows := s.mock.NewRows([]string{"is_exsists"})
	rows.AddRow(false)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Login, user.Email).
		WillReturnRows(rows)

	s.mock.ExpectExec("INSERT").
		WithArgs(user.Login, user.Email, user.Password,
			user.Name, s.repo.hostToSave+s.repo.defaultImagePath+s.repo.defaultAvatar).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Create(user)
	require.Nil(s.T(), err)
}
func (s *Suite) TestDeleteByLogin() {
	user := s.user

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Login).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.DeleteByLogin(user.Login)

	require.Nil(s.T(), err)

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Login).
		WillReturnResult(sqlmock.NewResult(50, 0))

	err = s.repo.DeleteByLogin(user.Login)

	require.Error(s.T(), err)

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Login).
		WillReturnError(entityerrors.DoesNotExist())

	err = s.repo.DeleteByLogin(user.Login)

	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}

func (s *Suite) TestDeleteById() {
	user := s.user

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.DeleteById(user.Id)

	require.Nil(s.T(), err)

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Id).
		WillReturnResult(sqlmock.NewResult(50, 0))

	err = s.repo.DeleteById(user.Id)

	require.Error(s.T(), err)

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Id).
		WillReturnError(entityerrors.DoesNotExist())

	err = s.repo.DeleteById(user.Id)

	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}
func (s *Suite) TestCheckPass() {
	user := s.user

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword[:])

	rows := s.mock.NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.ExpectQuery("SELECT").
		WithArgs(user.Login).
		WillReturnRows(rows)

	isCorrect, err := s.repo.CheckPass(user.Login, s.user.Password)
	require.True(s.T(), isCorrect)
	require.Nil(s.T(), err)

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Id).
		WillReturnResult(sqlmock.NewResult(50, 0))

	err = s.repo.DeleteById(user.Id)

	require.Error(s.T(), err)

	s.mock.ExpectExec("DELETE").
		WithArgs(user.Id).
		WillReturnError(entityerrors.DoesNotExist())

	err = s.repo.DeleteById(user.Id)

	require.Equal(s.T(), errors.Cause(err), entityerrors.DoesNotExist())
}

func (s *Suite) TestUpdateAvatarPath() {
	user := s.user

	rows := s.mock.NewRows([]string{"id", "login", "email", "password", "name", "avatar_path"})
	rows.AddRow(user.Id, user.Login, user.Email, user.Password, user.Name, user.Image)

	s.mock.ExpectExec("UPDATE").
		WithArgs(user.Id, user.Email, user.Name, s.repo.hostToSave+s.repo.defaultImagePath+user.Image, user.Password).
		WillReturnResult(sqlmock.NewResult(1, 1))

	 err := s.repo.UpdateAvatarPath(user, user.Image)

	require.Nil(s.T(), err)

}
