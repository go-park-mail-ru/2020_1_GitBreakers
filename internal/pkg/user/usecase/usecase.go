package usecase

import (
	"github.com/asaskevich/govalidator"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"mime/multipart"
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
func (UC *UCUser) UploadAvatar(User models.User, fileName *multipart.FileHeader, image multipart.File) error {
	byteImage, err := ioutil.ReadAll(image)
	if err != nil {
		return errors.Wrap(err, "err in uploadImage Usercase")
	}
	file, err := fileName.Open()
	if err != nil {
		return err
	}
	if err := checkFileContentType(file); err != nil {
		return err
	}

	if err := UC.RepUser.UploadAvatar(fileName.Filename, byteImage); err != nil {
		return errors.Wrap(err, "err in repo UploadAvatar")
	}

	if err := UC.RepUser.UpdateAvatarPath(User, fileName.Filename); err != nil {
		return errors.Wrap(err, "err in repo UpdateAvatarPath")
	}
	return nil
}

func checkFileContentType(file multipart.File) error {
	buffer := make([]byte, 512)

	_, err := file.Read(buffer)
	if err != nil || err == io.EOF {
		return err
	}

	contentType := http.DetectContentType(buffer)

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
