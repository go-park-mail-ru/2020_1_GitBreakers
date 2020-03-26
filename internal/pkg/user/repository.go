package user

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoUser interface {
	GetUserByIdWithPass(id int64) (models.User, error)
	GetUserByIdWithoutPass(id int64) (models.User, error)
	GetUserByLoginWithPass(login string) (models.User, error)
	GetByLoginWithoutPass(login string) (models.User, error)
	Create(newUser models.User) error
	Update(usrUpd models.User) error
	IsExists(user models.User) (bool, error)
	DeleteById(id int64) error
}
