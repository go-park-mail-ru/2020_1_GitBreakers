package grpc

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type NewsServer struct {
	UC codehub.UCCodeHubI
}

func NewUserServer(gserver *grpc.Server, UserUseCase codehub.UCCodeHubI) {
	//userServer := &NewsServer{UC: UserUseCase}
	//RegisterUserGrpcServer(gserver, userServer)
	reflection.Register(gserver)
}


