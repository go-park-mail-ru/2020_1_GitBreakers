package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UCUserWork struct {
	RepUser user.RepoUser
}

func (UC *UCUserWork) Create(user models.User) error {
	isExsist, err := UC.RepUser.IsExists(user)
	if err != nil {
		return errors.Wrap(err, "error in repo layer")
	}
	if isExsist {
		return errors.New("user with this login or email is already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//конвертим в строку
	user.Password = string(hashedPassword[:])
	if err := UC.RepUser.Create(user); err != nil {
		return errors.New("error with creating user")
	}
	return nil
}

func (UC *UCUserWork) Delete(user models.User) error {
	UC.RepUser.DeleteById(user.Id)
	return nil
}

func (UC *UCUserWork) Update(user models.User) error {
	isExsist, err := UC.RepUser.IsExists(user)
	if err != nil {
		return errors.Wrap(err, "error in repo layer")
	}
	if isExsist {
		return errors.New("user with this login or email is already exists")
	}
	//todo может быть перехеширование хеша пароля
	UC.RepUser.Update(user)
	return nil
}
func (UC *UCUserWork) GetByLogin(login string) (models.User, error) {
	return UC.RepUser.GetUserByLoginWithPass(login)
}
func (UC *UCUserWork) CheckPass(User models.User, pass string) (bool, error) {
	return UC.RepUser.CheckPass(User.Password, pass)
}
