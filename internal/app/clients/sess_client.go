package clients

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	session "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type SessClient struct {
	conn *grpc.ClientConn
}

func (c *SessClient) Connect() error {
	//todo port unhardcode
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	return nil
}
func (c *SessClient) CreateSess(UserID int64) (string, error) {
	client := session.NewSessionClient(c.conn)
	UserIDModel := session.UserID{UserID: int64(UserID)}
	sessIDModel, err := client.Create(context.Background(), &UserIDModel)
	if err != nil {
		return "", err
	}

	return sessIDModel.GetSessionID(), nil
}

func (c *SessClient) DelSess(SessID string) error {
	client := session.NewSessionClient(c.conn)
	SessIDModel := session.SessionID{SessionID: SessID}

	_, err := client.Delete(context.Background(), &SessIDModel)
	return err
}

func (c *SessClient) GetSess(SessID string) (models.Session, error) {
	client := session.NewSessionClient(c.conn)
	SessIDModel := session.SessionID{SessionID: SessID}

	SessModel, err := client.Get(context.Background(), &SessIDModel)
	if err != nil {
		return models.Session{}, nil
	}
	return models.Session{
		ID:     SessModel.GetID(),
		UserID: SessModel.GetUserID(),
	}, nil
}
