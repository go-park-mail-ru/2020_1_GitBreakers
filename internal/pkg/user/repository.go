package user

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoUser interface {
	Create(newUser models.User) error
	Update(userUpdates models.User) error
	GetProfileById(userId int64) (models.User, error)
	Delete(id int64) error
}
