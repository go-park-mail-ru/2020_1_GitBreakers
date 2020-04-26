package clients

import (
	"context"
	news "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/delivery/grpc"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type NewsClient struct {
	conn   *grpc.ClientConn
	client news.NewsClient
}

func NewNewsClient() (NewsClient, error) {
	newsClient := NewsClient{}
	if err := newsClient.Connect(); err != nil {
		return newsClient, err
	}
	return newsClient, nil
}

func (c *NewsClient) Connect() error {
	//todo port unhardcode
	conn, err := grpc.Dial(":8083", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	c.client = news.NewNewsClient(c.conn)
	return nil
}
func (c *NewsClient) GetNews(repoID, userID, limit, offset int64) (models.NewsSet, error) {
	req := &news.NewsReq{
		RepoID: repoID,
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
	newsResp, err := c.client.Get(context.Background(), req)
	if err != nil {
		return nil, err
	}
	newsList := make([]models.News, len(newsResp.News))
	for i, v := range newsResp.News {
		if err := copier.Copy(newsList[i], v); err != nil {
			return nil, err
		}
	}
	return newsList, err
}
