package session

import "time"

type SessRepo interface {
	Create(sid string, login string, expire time.Duration) (string, error)
	GetLoginById(sessionId string) (string, error)
	DeleteById(sid string) error
}
