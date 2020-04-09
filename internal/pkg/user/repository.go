package user

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
)

type RepoUser interface {
	GetUserByIdWithPass(id int) (models.User, error)
	GetUserByIdWithoutPass(id int) (models.User, error)
	GetUserByLoginWithPass(login string) (models.User, error)
	GetByLoginWithoutPass(login string) (models.User, error)
	GetLoginByID(id int) (string, error)
	GetIdByLogin(login string) (int, error)
	Create(newUser models.User) error
	Update(usrUpd models.User) error
	IsExists(user models.User) (bool, error)
	DeleteById(id int) error
	CheckPass(login string, newpass string) (bool, error)
	UploadAvatar(Name string, Content []byte) error
	UpdateAvatarPath(User models.User, Name string) error
	UserCanUpdate(user models.User) (bool, error)
}
