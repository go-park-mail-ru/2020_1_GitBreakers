package repository

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type DBWork struct {
	DB            *sqlx.DB
	DefaultAvatar string
}

func (repo DBWork) GetUserByIdWithPass(id int) (models.User, error) {
	User := models.User{}
	row := repo.DB.QueryRow("SELECT id, login, email, password,name,avatar_path  FROM users WHERE id = $1", id)

	err := row.Scan(&User.Id, &User.Login, &User.Email, &User.Password, &User.Name, &User.Image)
	if err != nil {
		return models.User{}, errors.Wrap(err, "error with db call")
	}
	return User, nil
}
func (repo DBWork) GetUserByIdWithoutPass(id int) (models.User, error) {
	storedUser, err := repo.GetUserByIdWithoutPass(id)
	if err != nil {
		return models.User{}, errors.Wrapf(err, "error in user GetById with id=%v", id)
	}
	storedUser.Password = ""
	return storedUser, nil
}

func (repo DBWork) GetUserByLoginWithPass(login string) (models.User, error) {
	storedUser := models.User{}
	row := repo.DB.QueryRow("SELECT id, login, email, password,name, avatar_path FROM users WHERE login = $1", login)

	err := row.Scan(&storedUser.Id, &storedUser.Login, &storedUser.Email, &storedUser.Password, &storedUser.Name, &storedUser.Image)

	if err != nil {
		return storedUser, errors.Wrapf(err, "error in user getUserByLogin with login=%v", login)
	}
	//todo можно возвращать более конкретную ошибку
	return storedUser, nil
}

func (repo DBWork) GetByLoginWithoutPass(login string) (models.User, error) {
	//todo может быть можно какой-либо метод сделать приватным
	storedUser, err := repo.GetUserByLoginWithPass(login)
	if err != nil {
		return models.User{}, errors.Wrapf(err, "error in user GetByLogin with login=%v", login)
	}
	storedUser.Password = ""
	return storedUser, nil
}

func (repo DBWork) Create(newUser models.User) error {
	if isExists, err := repo.IsExists(newUser); isExists || err != nil {
		if err != nil {
			newUser.Password = ""
			return errors.Wrapf(err, "error in user Create with newUser=%+v", newUser)
		}
		return entityerrors.AlreadyExist()
	}
	userQuery := `INSERT INTO users (login, email, password, name, avatar_path) VALUES ($1, $2, $3, $4,$5);`
	_, err := repo.DB.Exec(userQuery, newUser.Login, newUser.Email, newUser.Password, newUser.Name, repo.DefaultAvatar)

	if err != nil {
		return errors.Wrap(err, "error in user Create ")
	}

	return nil
}

//id юзера не меняется, достаточно скинуть в него новые данные
func (repo DBWork) Update(usrUpd models.User) error {
	result, err := repo.DB.Exec(
		"UPDATE users SET login = $2, email = $3,name=$4,avatar_path=$5,password=$6 WHERE id = $1",
		usrUpd.Id, usrUpd.Login, usrUpd.Email, usrUpd.Name, usrUpd.Image, usrUpd.Password)
	if err != nil {
		return errors.Wrap(err, "error with update data")
	}
	updLines, err := result.RowsAffected()
	if updLines <= 0 {
		//немного странный возврат ошибки
		return entityerrors.Invalid()
	}
	return nil
}

func (repo DBWork) IsExists(user models.User) (bool, error) {
	isExists := false
	err := repo.DB.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 OR login = $2 OR email = $3) as is_exists",
		user.Id, user.Login, user.Email).Scan(&isExists)
	if err != nil {
		//user.Password = ""
		return isExists, errors.Wrapf(err, "error in user IsExists with user=%+v", user)
	}
	return isExists, nil
}

func (repo DBWork) DeleteById(id int) error {
	result, err := repo.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return errors.Wrapf(err, "error in user DeleteById with id=%v", id)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "error in user DeleteById with id=%v", id)
	}

	if rowsAffected == 0 {
		return entityerrors.DoesNotExist()
	}

	return nil
}

//todo дублирование deleteByLogin and deleteById
func (repo DBWork) DeleteByLogin(login string) error {
	result, err := repo.DB.Exec("DELETE FROM users WHERE login = $1", login)
	if err != nil {
		return errors.Wrapf(err, "error in user DeleteByLogin with login=%v", login)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "error in user DeleteByLogin with login=%v", login)
	}

	if rowsAffected == 0 {
		return entityerrors.DoesNotExist()
	}

	return nil
}
func (repo DBWork) CheckPass(oldpass string, newpass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(oldpass), []byte(newpass))
	if err != nil {
		return false, err
	}
	return true, nil
}
