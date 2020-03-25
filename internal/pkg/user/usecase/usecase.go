package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/repository"
)

type UCUserWork struct {
	RepUser user.RepoUser
}

func (UC *UCUserWork) Create(user models.User) error {
	//call RepUser.create()
	return nil
}

func (UC *UCUserWork) Delete(user models.User) error {
	//call RepUser.create()
	return nil
}

func (UC *UCUserWork) Update(user models.User) error {
	//call RepUser.create()
	return nil
}
func (UC *UCUserWork) GetByLogin(userid string) error {
	//call RepUser.create()
	return nil
}
