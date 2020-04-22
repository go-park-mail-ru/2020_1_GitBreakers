package session

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"time"
)

type SessServer struct {
	UC session.UCSession
}

func NewSessServer(gserver *grpc.Server, SessionUcase session.UCSession) {
	sessServer := &SessServer{
		UC: SessionUcase,
	}
	RegisterSessionServer(gserver, sessServer)
	reflection.Register(gserver)
}
func (h *SessServer) Create(ctx context.Context, req *UserID) (*SessionID, error) {
	sess := models.Session{UserID: int(req.GetUserID())}
	//todo duration in hardcode
	sessID, err := h.UC.Create(sess, 48*time.Hour)
	if err != nil {
		return &SessionID{}, err
	}
	return &SessionID{SessionID: sessID}, nil
}
func (h *SessServer) Delete(ctx context.Context, req *SessionID) (*empty.Empty, error) {
	err := h.UC.Delete(req.GetSessionID())
	if err != nil {
		return &empty.Empty{}, err
	}
	return &empty.Empty{}, nil
}
func (h *SessServer) Get(ctx context.Context, req *SessionID) (*SessionModel, error) {
	sess, err := h.UC.GetByID(req.GetSessionID())
	if err != nil {
		return &SessionModel{
			ID:     "",
			UserID: -1,
		}, err
	}

	return &SessionModel{
		ID:     sess.ID,
		UserID: int64(sess.UserID),
	}, nil
}
