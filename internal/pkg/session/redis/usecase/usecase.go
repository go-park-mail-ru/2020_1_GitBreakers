package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"time"
)

type SessionUC struct {
	RepoSession session.SessRepo //содержит в себе класс репо+его методы
}

func (sessUC *SessionUC) Create(session models.Session, expire time.Duration) (string, error) {
	return sessUC.RepoSession.Create(session, expire)
}

func (sessUC *SessionUC) Delete(sessionID string) error {
	return sessUC.RepoSession.DeleteById(sessionID)

}

func (sessUC *SessionUC) GetByID(sid string) (models.Session, error) {
	return sessUC.RepoSession.GetSessById(sid)
}
