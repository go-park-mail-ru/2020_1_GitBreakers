package user

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"mime/multipart"
)

type UCUser interface {
	Create(user models.User) error
	Delete(user models.User) error
	Update(userID int, user models.User) error
	GetByLogin(login string) (models.User, error)
	GetByID(userID int) (models.User, error)
	CheckPass(login string, pass string) (bool, error)
	//todo mime type to delivery layer
	UploadAvatar(User models.User, fileName *multipart.FileHeader, file multipart.File) error
}
