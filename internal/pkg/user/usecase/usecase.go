package usecase

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UCUser struct {
	RepUser user.RepoUser
}

func (UC *UCUser) Create(user models.User) error {
	isExsist, err := UC.RepUser.IsExists(user)
	if err != nil {
		return errors.Wrap(err, "error in repo layer")
	}
	if isExsist {
		return entityerrors.AlreadyExist()
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	//конвертим в строку
	user.Password = string(hashedPassword[:])
	if err := UC.RepUser.Create(user); err != nil {
		return errors.New("error with creating user")
	}
	return nil
}

func (UC *UCUser) Delete(user models.User) error {
	if err := UC.RepUser.DeleteByID(user.ID); err != nil {
		return errors.Wrap(err, "can't delete user")
	}
	return nil
}

func (UC *UCUser) Update(userid int, newUserData models.User) error {
	oldUserData, err := UC.RepUser.GetUserByIDWithPass(userid)
	if err != nil {
		return errors.Wrap(err, "error in repo layer")
	}
	if govalidator.IsByteLength(newUserData.Name, 5, 128) {
		oldUserData.Name = newUserData.Name
	}
	if govalidator.IsEmail(newUserData.Email) {
		oldUserData.Email = newUserData.Email
	}
	if govalidator.IsByteLength(newUserData.Password, 5, 128) {
		pass, err := bcrypt.GenerateFromPassword([]byte(newUserData.Password), bcrypt.DefaultCost)
		if err != nil {
			return errors.Wrap(err, "error in bcrypt")
		}
		oldUserData.Password = string(pass[:])
	}
	if canUpdate, err := UC.RepUser.UserCanUpdate(oldUserData); !canUpdate || err != nil {
		if !canUpdate {
			return entityerrors.AlreadyExist()
		}
		return err
	}
	if err := UC.RepUser.Update(oldUserData); err != nil {
		return err
	}
	return nil
}

func (UC *UCUser) GetByLogin(login string) (models.User, error) {
	return UC.RepUser.GetByLoginWithoutPass(login)
}
func (UC *UCUser) GetByID(userId int) (models.User, error) {
	return UC.RepUser.GetUserByIDWithoutPass(userId)
}
func (UC *UCUser) CheckPass(login string, pass string) (bool, error) {
	return UC.RepUser.CheckPass(login, pass)
}
func (UC *UCUser) UploadAvatar(UserID int, fileName string, fileData []byte) error {

	if err := checkFileContentType(fileData); err != nil {
		return err
	}

	if err := UC.RepUser.UploadAvatar(fileName, fileData); err != nil {
		return errors.Wrap(err, "err in repo UploadAvatar")
	}
	UserModel, err := UC.RepUser.GetUserByIDWithPass(UserID)

	if err != nil {
		return errors.Wrap(err, "err in GetUserByLoginWithPass")
	}

	if err := UC.RepUser.UpdateAvatarPath(UserModel, fileName); err != nil {
		return errors.Wrap(err, "err in repo UpdateAvatarPath")
	}
	return nil
}

func checkFileContentType(fileContent []byte) error {

	contentType := http.DetectContentType(fileContent)

	for _, r := range allowedContentType {
		if contentType == r {
			return nil
		}
	}
	return errors.New("this content type is not allowed")
}

var allowedContentType = []string{
	"image/png",
	"image/jpeg",
}
