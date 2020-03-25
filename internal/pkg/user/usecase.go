package user

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type UCUser interface {
	GetById(id int64) (models.User, error)
	GetByLogin(login string) (models.User, error)
	Create(user models.User) error
	IsExists(user models.User) (bool, error)
}
