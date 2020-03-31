package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	uuid "github.com/satori/go.uuid"
	"time"
)

type SessionUC struct {
	RepoSession session.SessRepo //содержит в себе класс репо+его методы
}

func (sessUC *SessionUC) Create(user models.User, expires time.Duration) (string, error) {
	sid := uuid.NewV4().String()
	_, err := sessUC.RepoSession.Create(sid, user.Login, expires)
	if err != nil {
		return "", err
	}
	return sid, nil
}

func (sessUC *SessionUC) Delete(sessionID string) error {
	err := sessUC.RepoSession.DeleteById(sessionID)
	if err != nil {
		return errors.New("error with delete session")
	}
	return nil
}

func (sessUC *SessionUC) GetLoginBySessID(sid string) (string, error) {
	login, err := sessUC.RepoSession.GetLoginById(sid)
	if err != nil {
		return "", nil
	}
	return login, nil
}
