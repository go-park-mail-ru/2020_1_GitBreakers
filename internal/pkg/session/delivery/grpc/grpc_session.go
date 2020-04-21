package session

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"time"
)

type GRPCServer struct {
	UC session.UCSession
}

func (h *GRPCServer) Create(ctx context.Context, req *CreateReq) (*CreateResp, error) {
	sess := models.Session{UserID: int(req.GetUserID())}
	//todo duration in hardcode
	sessID, err := h.UC.Create(sess, 48*time.Hour)
	if err != nil {
		return &CreateResp{}, err
	}
	return &CreateResp{SessionID: sessID}, nil
}

