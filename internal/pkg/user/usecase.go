package user

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
)

type UCUser interface {
	Create(user models.User) error
	Delete(user models.User) error
	Update(userID int64, user models.User) error
	GetByLogin(login string) (models.User, error)
	GetByID(userID int64) (models.User, error)
	CheckPass(login string, pass string) (bool, error)
	UploadAvatar(UserID int64, fileName string, fileData []byte) error
}
