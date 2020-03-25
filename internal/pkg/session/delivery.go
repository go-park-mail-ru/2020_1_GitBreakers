package session

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"net/http"
)

type SessDelivery interface {
	Create(user *models.User) (*http.Cookie, error)
	Delete(user models.User) error
	GetLoginBySessionID(sessionID string) (string, error)
}
