package news

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"github.com/jinzhu/copier"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type NewsServerOwn struct {
	UC codehub.UCCodeHubI
}

func NewNewsServer(gserver *grpc.Server, UserUseCase codehub.UCCodeHubI) {
	newsServer := &NewsServerOwn{UC: UserUseCase}
	RegisterNewsServer(gserver, newsServer)
	reflection.Register(gserver)
}

func (s NewsServerOwn) Get(ctx context.Context, in *NewsReq) (*NewsResp, error) {
	newsList, err := s.UC.GetNews(in.GetRepoID(), in.GetUserID(), in.GetLimit(), in.GetOffset())
	if err != nil {
		return &NewsResp{}, err
	}
	//todo debug может копировать хренотень
	newsListProto := make([]*NewsModel, len(newsList))
	for i, v := range newsList {
		if err := copier.Copy(newsListProto[i], v); err != nil {
			return &NewsResp{}, err
		}
	}

	return &NewsResp{News: newsListProto}, nil
}
