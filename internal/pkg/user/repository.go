package user

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
)

type RepoUser interface {
	GetUserByIDWithPass(ID int64) (models.User, error)
	GetUserByIDWithoutPass(ID int64) (models.User, error)
	GetUserByLoginWithPass(login string) (models.User, error)
	GetByLoginWithoutPass(login string) (models.User, error)
	GetLoginByID(ID int64) (string, error)
	GetIDByLogin(login string) (int64, error)
	Create(newUser models.User) error
	Update(usrUpd models.User) error
	IsExists(user models.User) (bool, error)
	DeleteByID(ID int64) error
	CheckPass(login string, newpass string) (bool, error)
	UploadAvatar(Name string, Content []byte) error
	UpdateAvatarPath(User models.User, Name string) error
	UserCanUpdate(user models.User) (bool, error)
	GetByEmail(email string) (models.User, error)
}
