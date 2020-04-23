package usergrpc

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
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
	User := models.User{
		ID:       int(in.GetId()),
		Password: in.GetPassword(),
		Name:     in.GetName(),
		Login:    in.GetLogin(),
		Image:    in.GetImage(),
		Email:    in.GetEmail(),
	}
	return nil, h.UC.Create(User)
}
func (h *UserServer) GetByLogin(ctx context.Context, in *LoginModel) (*UserModel, error) {
	User, err := h.UC.GetByLogin(in.GetLogin())
	userModel := h.transfToUserModelGRPC(&User)

	return &userModel, err
}
func (h *UserServer) UpdateUser(ctx context.Context, in *UserUpdateModel) (*empty.Empty, error) {
	UserID := int(in.GetUserID())
	User := h.transfToModelsUser(in.GetUserData())
	return nil, h.UC.Update(UserID, User)
}
func (h *UserServer) CheckPass(ctx context.Context, in *CheckPassModel) (*CheckPassResp, error) {
	isCorrect, err := h.UC.CheckPass(in.GetLogin(), in.GetPass())
	return &CheckPassResp{IsCorrect: isCorrect}, err
}
func (h *UserServer) GetByID(ctx context.Context, in *UserIDModel) (*UserModel, error) {
	User, err := h.UC.GetByID(int(in.GetUserID()))
	userModel := h.transfToUserModelGRPC(&User)
	return &userModel, err
}
func (h *UserServer) UploadAvatar(UserGrpc_UploadAvatarServer) error {
	return nil
}

//функции вспомогательные для преобразования from/to models.User UserModel(GRPC)
func (h *UserServer) transfToModelsUser(in *UserModel) models.User {
	return models.User{
		ID:       int(in.GetId()),
		Password: in.GetPassword(),
		Name:     in.GetName(),
		Login:    in.GetLogin(),
		Image:    in.GetImage(),
		Email:    in.GetEmail(),
	}
}

//функции вспомогательные для преобразования from/to models.User UserModel(GRPC)
func (h *UserServer) transfToUserModelGRPC(in *models.User) UserModel {
	return UserModel{
		Id:       int64(in.ID),
		Password: in.Password,
		Name:     in.Name,
		Login:    in.Login,
		Image:    in.Image,
		Email:    in.Email,
	}
}
