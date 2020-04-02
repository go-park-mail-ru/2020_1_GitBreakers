package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"net/http"
	"time"
)

type SessionHttp struct {
	SessUC     session.UCSession
	ExpireTime time.Duration
}

func (UC *SessionHttp) Create(userID int) (http.Cookie, error) {
	baseSess := models.Session{UserId: userID}
	cretedSess, err := UC.SessUC.Create(baseSess, UC.ExpireTime)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    cretedSess,
		HttpOnly: true,
		Expires:  time.Now().Add(UC.ExpireTime),
		Path:     "/",
	}, nil
}

func (UC *SessionHttp) Delete(sessID string) error {
	return UC.SessUC.Delete(sessID)
}

func (UC *SessionHttp) GetBySessID(sessionID string) (models.Session, error) {
	return UC.SessUC.GetByID(sessionID)
}
