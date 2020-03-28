package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"net/http"
	"time"
)

type SessionHttpWork struct {
	SessUC     session.UCSession
	ExpireTime time.Duration
}

func (UC *SessionHttpWork) Create(user models.User) (http.Cookie, error) {
	sid, err := UC.SessUC.Create(user, UC.ExpireTime)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    sid,
		HttpOnly: true,
		Expires:  time.Now().Add(UC.ExpireTime),
	}, nil
}

func (UC *SessionHttpWork) Delete(sessID string) error {
	return UC.SessUC.Delete(sessID)
}

func (UC *SessionHttpWork) GetLoginBySessID(sessionID string) (string, error) {
	return UC.SessUC.GetLoginBySessID(sessionID)
}
