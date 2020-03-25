package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/repository"
)

type SessionUCWork struct {
	repoSession session.SessRepo //содержит в себе класс репо+его методы
}

func (sessUC *SessionUCWork) Create(user models.User) error {
	//sessUC.repoSession.Create(user.Login,time.ParseDuration("10m"))
	return nil
}

func (sessUC *SessionUCWork) Delete(user models.User) error {
	//sessUC.repoSession.Delete(user.Login,time.ParseDuration("10m"))
	return nil
}

func (sessUC *SessionUCWork) GetLoginBySessId(id string) (string, error) {
	//sessUC.repoSession.find(user.Login,time.ParseDuration("10m"))
	return "result", nil
}
