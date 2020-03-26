package session

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"time"
)

type SessDelivery interface {
	Create(user models.User, expires time.Duration) (string, error)
	Delete(sessionID string) error
	GetLoginBySessID(sessionID string) (string, error)
}
