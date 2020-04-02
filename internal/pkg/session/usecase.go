package session

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"time"
)

type UCSession interface {
	Create(session models.Session, expire time.Duration) (string, error)
	Delete(sessionID string) error
	GetByID(sid string) (models.Session, error)
}
