package session

import (
	"time"
)

type SessRepo interface {
	Create(id string, login string, expire time.Duration) error
	GetLogin(sessionId string) error
	Delete(sessionId string) error
}
