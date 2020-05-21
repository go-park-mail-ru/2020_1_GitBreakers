package clients

import (
	"context"
	news "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/delivery/grpc"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/golang/protobuf/ptypes"
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
	c.client = news.NewNewsClient(conn)
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
	newsList := make(models.NewsSet, 0)
	for _, v := range newsResp.News {
		var newsModel models.News

		newsModel.ID = v.ID
		newsModel.AuthorID = v.AuthorID
		newsModel.Date, err = ptypes.Timestamp(v.Date)
		newsModel.Mess = v.Message
		newsModel.RepoID = v.RepoID
		newsModel.AuthorLogin = v.AuthorLogin
		newsModel.AuthorImage = v.AuthorImage

		newsList = append(newsList, newsModel)
	}

	return newsList, err
}
