package user

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type RepoUser interface {
	Create(newUser models.User) error
	Update(userUpdates models.User) error
	SaveAvatar(user *models.User, filepath string) error
}
