package session

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"time"
)

type UCSession interface {
	Create(user models.User, expires time.Duration) (string, error)
	Delete(sessionID string) error
	GetLoginBySessionID(sessionID string) (string, error)
}
