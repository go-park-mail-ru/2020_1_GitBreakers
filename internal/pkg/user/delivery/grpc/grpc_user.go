package usergrpc

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type UserServer struct {
	UC user.UCUser
}

func NewUserServer(gserver *grpc.Server, UserUseCase user.UCUser) {
	userServer := &UserServer{UC: UserUseCase}
	RegisterUserGrpcServer(gserver, userServer)
	reflection.Register(gserver)
}
func (h *UserServer) Create(ctx context.Context, in *UserModel) (*empty.Empty, error) {
	return nil, nil
}
func (h *UserServer) GetByLogin(ctx context.Context, in *LoginModel) (*UserModel, error) {
	return nil, nil
}
func (h *UserServer) UpdateUser(ctx context.Context, in *UserUpdateModel) (*empty.Empty, error) {
	return nil, nil
}
func (h *UserServer) CheckPass(ctx context.Context, in *CheckPassModel) (*CheckPassResp, error) {
	return nil, nil
}
func (h *UserServer) GetByID(ctx context.Context, in *UserIDModel) (*UserModel, error) {
	return nil, nil
}
//todo не хочет удовлетворять интерфейсу
func (h *UserServer) UploadAvatar(ctx context.Context, in *UserUpdateModel) error {
	return nil
}
