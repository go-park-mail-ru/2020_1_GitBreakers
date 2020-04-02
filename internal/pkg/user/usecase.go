package user

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"mime/multipart"
)

type UCUser interface {
	Create(user models.User) error
	Delete(user models.User) error
	Update(userid int, user models.User) error
	GetByLogin(login string) (models.User, error)
	GetByID(userId int) (models.User, error)
	CheckPass(User models.User, pass string) (bool, error)
	UploadAvatar(User models.User, fileName *multipart.FileHeader, file multipart.File) error
}
