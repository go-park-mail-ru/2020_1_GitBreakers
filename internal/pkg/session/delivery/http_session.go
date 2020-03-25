package delivery

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"net/http"
	"time"
)

type SessionHttpWork struct {
	SessUC session.UCSession
}

func (UC *SessionHttpWork) Create(user *models.User) (*http.Cookie, error) {
	//create session
	timeCookie, _ := time.ParseDuration("1h")
	return &http.Cookie{
		Name:     "session_id",
		Value:    "somesessionhash",
		HttpOnly: true,
		Expires:  time.Now().Add(timeCookie),
	}, nil
}

func (UC *SessionHttpWork) Delete(user models.User) error {
	return UC.SessUC.Delete(user)
}

func (UC *SessionHttpWork) GetLoginBySessionID(sessionID string) (string, error) {
	return UC.SessUC.GetLoginBySessId(sessionID)
}
