package interfaces

import "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"

type SessClientI interface {
	CreateSess(UserID int64) (string, error)
	DelSess(SessID string) error
	GetSess(SessID string) (models.Session, error)
}
