package session

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"time"
)

type SessRepo interface {
	Create(session models.Session, expire time.Duration) (string, error)
	GetSessByID(sessionID string) (models.Session, error)
	DeleteByID(sessionID string) error
}
