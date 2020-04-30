package interfaces

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type UserClientI interface {
	Create(User models.User) error
	Update(userID int64, newUserData models.User) error
	GetByLogin(login string) (models.User, error)
	GetByID(userID int64) (models.User, error)
	CheckPass(login string, pass string) (bool, error)
	UploadAvatar(UserID int64, fileName string, fileData []byte, fileSize int64) error
}
