package user

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type UCUser interface {
	Create(user models.User) error
	Delete(user models.User) error
	Update(user models.User) error
	GetByLogin(login string) (models.User, error)
	CheckPass(User models.User, pass string) (bool, error)
}
