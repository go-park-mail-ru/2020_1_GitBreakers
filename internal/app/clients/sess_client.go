package clients

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	session "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/grpc"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type SessClient struct {
	conn   *grpc.ClientConn
	client session.SessionClient
}

func NewSessClient() (SessClient, error) {
	sessClient := SessClient{}
	if err := sessClient.Connect(); err != nil {
		return sessClient, err
	}
	return sessClient, nil
}

func (c *SessClient) Connect() error {
	//todo port unhardcode
	conn, err := grpc.Dial(":8081", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	c.client = session.NewSessionClient(c.conn)
	return nil
}
func (c *SessClient) CreateSess(UserID int64) (string, error) {

	UserIDModel := session.UserID{UserID: int64(UserID)}
	sessIDModel, err := c.client.Create(context.Background(), &UserIDModel)
	if err != nil {
		return "", err
	}

	return sessIDModel.GetSessionID(), nil
}

func (c *SessClient) DelSess(SessID string) error {
	SessIDModel := session.SessionID{SessionID: SessID}

	_, err := c.client.Delete(context.Background(), &SessIDModel)
	return err
}

func (c *SessClient) GetSess(SessID string) (models.Session, error) {
	SessIDModel := session.SessionID{SessionID: SessID}

	SessModel, err := c.client.Get(context.Background(), &SessIDModel)
	if err != nil {
		return models.Session{}, err
	}
	return models.Session{
		ID:     SessModel.GetID(),
		UserID: SessModel.GetUserID(),
	}, nil
}
