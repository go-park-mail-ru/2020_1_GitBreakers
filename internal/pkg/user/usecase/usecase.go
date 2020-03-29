package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
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
func (UC *UCUserWork) UploadAvatar(User models.User, fileName *multipart.FileHeader, image multipart.File) error {
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
