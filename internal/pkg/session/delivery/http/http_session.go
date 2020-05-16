package http

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"net/http"
	"time"
)

type SessionHttp struct {
	ExpireTime time.Duration
	Client     interfaces.SessClientI
}

func (UC *SessionHttp) Create(userID int64) (http.Cookie, error) {
	cretedSess, err := UC.Client.CreateSess(userID)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     "session_id",
		Value:    cretedSess,
		HttpOnly: true,
		Expires:  time.Now().Add(UC.ExpireTime),
		Path:     "/",
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}, nil
}

func (UC *SessionHttp) Delete(sessID string) error {
	return UC.Client.DelSess(sessID)
}

func (UC *SessionHttp) GetBySessID(sessionID string) (models.Session, error) {
	return UC.Client.GetSess(sessionID)
}
