package http

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/interfaces"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"net/http"
	"time"
)

type SessionHttp struct {
	CookieName       string
	CookieExpireTime time.Duration
	CookieSecure     bool
	CookieSiteMode   http.SameSite
	CookiePath       string
	Client           interfaces.SessClientI
}

func (UC *SessionHttp) Create(userID int64) (http.Cookie, error) {
	cretedSess, err := UC.Client.CreateSess(userID)
	if err != nil {
		return http.Cookie{}, err
	}
	return http.Cookie{
		Name:     UC.CookieName,
		Value:    cretedSess,
		HttpOnly: true,
		Expires:  time.Now().Add(UC.CookieExpireTime),
		Path:     UC.CookiePath,
		SameSite: UC.CookieSiteMode,
		Secure:   UC.CookieSecure,
	}, nil
}

func (UC *SessionHttp) Delete(sessID string) error {
	return UC.Client.DelSess(sessID)
}

func (UC *SessionHttp) GetBySessID(sessionID string) (models.Session, error) {
	return UC.Client.GetSess(sessionID)
}
