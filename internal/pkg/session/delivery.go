package session

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"net/http"
)

type SessDelivery interface {
	Create(userID int) (http.Cookie, error)
	Delete(sessID string) error
	GetBySessID(sessionID string) (models.Session, error)
}
