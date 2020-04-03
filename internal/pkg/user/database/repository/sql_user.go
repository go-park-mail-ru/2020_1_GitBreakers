package repository

import (
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
func (repo UserRepo) GetUserByIdWithPass(id int) (models.User, error) {
	User := models.User{}
	row := repo.db.QueryRow("SELECT id, login, email, password,name,avatar_path  FROM users WHERE id = $1", id)

	err := row.Scan(&User.Id, &User.Login, &User.Email, &User.Password, &User.Name, &User.Image)
	if err != nil {
		return models.User{}, errors.Wrap(err, "error with db call")
	}
	return User, nil
}
func (repo UserRepo) GetUserByIdWithoutPass(id int) (models.User, error) {
	storedUser, err := repo.GetUserByIdWithPass(id)
	if err != nil {
		return models.User{}, errors.Wrapf(err, "error in user GetById with id=%v", id)
	}
	storedUser.Password = ""
	return storedUser, nil
}

func (repo UserRepo) GetUserByLoginWithPass(login string) (models.User, error) {
	storedUser := models.User{}
	row := repo.db.QueryRow("SELECT id, login, email, password,name, avatar_path FROM users WHERE login = $1", login)

	err := row.Scan(&storedUser.Id, &storedUser.Login, &storedUser.Email, &storedUser.Password, &storedUser.Name, &storedUser.Image)

	if err != nil {
		return storedUser, errors.Wrapf(err, "error in user getUserByLogin with login=%v", login)
	}
	//todo можно возвращать более конкретную ошибку
	return storedUser, nil
}

func (repo UserRepo) GetByLoginWithoutPass(login string) (models.User, error) {
	//todo может быть можно какой-либо метод сделать приватным
	storedUser, err := repo.GetUserByLoginWithPass(login)
	if err != nil {
		return models.User{}, errors.Wrapf(err, "error in user GetByLogin with login=%v", login)
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

func (repo UserRepo) IsExists(user models.User) (bool, error) {
	isExists := false
	err := repo.db.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM users WHERE login = $1 OR email = $2) as is_exists",
		user.Login, user.Email).Scan(&isExists)
	if err != nil {
		//user.Password = ""
		return isExists, errors.Wrapf(err, "error in user IsExists with user=%+v", user)
	}
	return isExists, nil
}
func (repo UserRepo) UserCanUpdate(user models.User) (bool, error) {
	usercount := 0
	err := repo.db.Get(&usercount,
		`select count(*) from users where login = $1 OR email = $2`, user.Login, user.Email)
	if err != nil {
		return false, err
	}
	if usercount > 1 {
		return false, nil
	} else {
		return true, nil
	}
}

func (repo UserRepo) DeleteById(id int) error {
	result, err := repo.db.Exec("DELETE FROM users WHERE id = $1", id)
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
func (repo UserRepo) CheckPass(oldpass string, newpass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(oldpass), []byte(newpass))
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
