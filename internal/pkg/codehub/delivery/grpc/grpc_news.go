package news

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
)

type NewsServerOwn struct {
	UC codehub.UCCodeHubI
}

func NewNewsServer(gserver *grpc.Server, CodehubUseCase codehub.UCCodeHubI) {
	newsServer := &NewsServerOwn{UC: CodehubUseCase}
	RegisterNewsServer(gserver, newsServer)
}

func (s NewsServerOwn) Get(ctx context.Context, in *NewsReq) (*NewsResp, error) {
	newsList, err := s.UC.GetNews(in.GetRepoID(), in.GetUserID(), in.GetLimit(), in.GetOffset())
	if err != nil {
		return &NewsResp{}, err
	}

	newsListProto := make([]*NewsModel, 0)
	for _, v := range newsList {
		var temp NewsModel

		temp.ID = v.ID
		temp.AuthorID = v.AuthorID
		temp.Date, _ = ptypes.TimestampProto(v.Date)
		temp.RepoID = v.RepoID
		temp.Message = v.Mess
		temp.AuthorLogin = v.AuthorLogin
		temp.AuthorImage = v.AuthorImage

		newsListProto = append(newsListProto, &temp)
	}

	return &NewsResp{News: newsListProto}, nil
}
