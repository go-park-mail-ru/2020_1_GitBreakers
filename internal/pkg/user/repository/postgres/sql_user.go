package postgres

import (
	"database/sql"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
)

type UserRepo struct {
	db               *sqlx.DB
	defaultAvatar    string
	defaultImagePath string
	hostToSave       string
}

func NewUserRepo(conn *sqlx.DB, defAva string, defPath string, defHost string) UserRepo {
	return UserRepo{
		db:               conn,
		defaultAvatar:    defAva,
		defaultImagePath: defPath,
		hostToSave:       defHost,
	}
}
func (repo UserRepo) GetUserByIDWithPass(id int) (models.User, error) {
	User := models.User{}
	err := repo.db.Get(&User, "SELECT id, login, email, password,name,avatar_path  FROM users WHERE id = $1", id)

	switch {
	case err == sql.ErrNoRows:
		return User, entityerrors.DoesNotExist()
	case err != nil:
		return User, errors.Wrap(err, "error while scanning in user")
	}

	return User, nil
}
func (repo UserRepo) GetUserByIDWithoutPass(id int) (models.User, error) {
	storedUser, err := repo.GetUserByIDWithPass(id)
	if err != nil {
		return models.User{}, err
	}
	storedUser.Password = ""
	return storedUser, nil
}

func (repo UserRepo) GetUserByLoginWithPass(login string) (models.User, error) {
	User := models.User{}
	err := repo.db.Get(&User, "SELECT id, login, email, password,name, avatar_path FROM users WHERE login = $1", login)
	switch {
	case err == sql.ErrNoRows:
		return User, entityerrors.DoesNotExist()
	case err != nil:
		return User, errors.Wrap(err, "error while scanning in repository")
	}
	return User, nil
}
func (repo UserRepo) GetByLoginWithoutPass(login string) (models.User, error) {
	storedUser, err := repo.GetUserByLoginWithPass(login)
	if err != nil {
		return models.User{}, err
	}
	storedUser.Password = ""
	return storedUser, nil
}

func (repo UserRepo) Create(newUser models.User) error {
	if isExists, err := repo.IsExists(newUser); isExists || err != nil {
		if err != nil {
			newUser.Password = ""
			return errors.Wrapf(err, "error in user Create with newUser=%+v", newUser)
		}
		return entityerrors.AlreadyExist()
	}
	userQuery := `INSERT INTO users (login, email, password, name, avatar_path) VALUES ($1, $2, $3, $4,$5);`
	_, err := repo.db.Exec(userQuery, newUser.Login, newUser.Email, newUser.Password,
		newUser.Name, repo.hostToSave+repo.defaultImagePath+repo.defaultAvatar)

	if err != nil {
		return errors.Wrap(err, "error in user Create ")
	}

	return nil
}

//id юзера не меняется, достаточно скинуть в него новые данные
func (repo UserRepo) Update(usrUpd models.User) error {
	result, err := repo.db.Exec(
		"UPDATE users SET  email = $2,name=$3,avatar_path=$4,password=$5 WHERE id = $1",
		usrUpd.ID, usrUpd.Email, usrUpd.Name, usrUpd.Image, usrUpd.Password)
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

func (repo UserRepo) IsExists(user models.User) (bool, error) {
	isExists := false
	err := repo.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 OR email = $2) as is_exists",
		user.Login, user.Email).Scan(&isExists)
	if err != nil {
		return isExists, errors.Wrapf(err, "error in user IsExists with user=%+v", user)
	}
	return isExists, nil
}
func (repo UserRepo) UserCanUpdate(user models.User) (bool, error) {
	usercount := 0
	err := repo.db.QueryRow(
		`SELECT count(*) from users where login = $1 OR email = $2`, user.Login, user.Email).Scan(&usercount)
	if err != nil {
		return false, err
	}
	if usercount > 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo UserRepo) DeleteByID(id int) error {
	result, err := repo.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return errors.Wrapf(err, "error in user DeleteByID with id=%v", id)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "error in user DeleteByID with id=%v", id)
	}

	if rowsAffected == 0 {
		return entityerrors.DoesNotExist()
	}

	return nil
}

func (repo UserRepo) DeleteByLogin(login string) error {
	result, err := repo.db.Exec("DELETE FROM users WHERE login = $1", login)
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
func (repo UserRepo) CheckPass(login string, newpass string) (bool, error) {
	User, err := repo.GetUserByLoginWithPass(login)
	if err != nil {
		return false, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(User.Password), []byte(newpass))
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo UserRepo) UpdateAvatarPath(User models.User, Name string) error {
	User.Image = repo.hostToSave + repo.defaultImagePath + Name
	if err := repo.Update(User); err != nil {
		return errors.Wrap(err, "error in db")
	}
	return nil
}
func (repo UserRepo) UploadAvatar(Name string, Content []byte) error {
	if err := ioutil.WriteFile(`.`+repo.defaultImagePath+Name, Content, 0644); err != nil {
		return errors.Wrap(err, " in repo user upload avatar")
	}
	return nil
}

func (repo UserRepo) GetLoginByID(id int) (string, error) {
	var login string
	err := repo.db.QueryRow("SELECT login FROM users WHERE id = $1", id).Scan(&login)
	switch {
	case err == sql.ErrNoRows:
		return login, entityerrors.DoesNotExist()
	case err != nil:
		return login, errors.Wrap(err, "error in user GetLoginByID")
	}
	return login, err
}

func (repo UserRepo) GetIDByLogin(login string) (int, error) {
	var id int
	err := repo.db.QueryRow("SELECT id FROM users WHERE login = $1", login).Scan(&id)
	switch {
	case err == sql.ErrNoRows:
		return id, entityerrors.DoesNotExist()
	case err != nil:
		return id, errors.Wrap(err, "error in user GetLoginByID")
	}
	return id, nil
}