package usergrpc

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
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
		ID:       in.GetID(),
		Password: in.GetPassword(),
		Name:     in.GetName(),
		Login:    in.GetLogin(),
		Image:    in.GetImage(),
		Email:    in.GetEmail(),
	}
	return &empty.Empty{}, h.UC.Create(User)
}
func (h *UserServer) GetByLogin(ctx context.Context, in *LoginModel) (*UserModel, error) {
	User, err := h.UC.GetByLogin(in.GetLogin())
	userModel := h.transfToUserModelGRPC(&User)

	return &userModel, err
}
func (h *UserServer) UpdateUser(ctx context.Context, in *UserUpdateModel) (*empty.Empty, error) {
	UserID := in.GetUserID()
	User := h.transfToModelsUser(in.GetUserData())
	return &empty.Empty{}, h.UC.Update(UserID, User)
}
func (h *UserServer) CheckPass(ctx context.Context, in *CheckPassModel) (*CheckPassResp, error) {
	isCorrect, err := h.UC.CheckPass(in.GetLogin(), in.GetPass())
	return &CheckPassResp{IsCorrect: isCorrect}, err
}
func (h *UserServer) GetByID(ctx context.Context, in *UserIDModel) (*UserModel, error) {
	User, err := h.UC.GetByID(in.GetUserID())
	userModel := h.transfToUserModelGRPC(&User)
	return &userModel, err
}
func (h *UserServer) UploadAvatar(stream UserGrpc_UploadAvatarServer) error {
	buf := make([]byte, 0)
	userAvatar := UserAvatarModel{}
	for {
		tempModel, err := stream.Recv()
		switch {
		case err == io.EOF:
			buf = append(buf, tempModel.GetChunk()...)
			goto END
		case err != nil:
			err = errors.Wrapf(err, "failed unexpectadely while reading chunks from stream")
			return err
		default:
			//todo каждый раз переопределяются overhead
			userAvatar.FileName = tempModel.FileName
			userAvatar.UserID = tempModel.UserID
			buf = append(buf, tempModel.GetChunk()...)
		}
	}
END:
	// once the transmission finished, send the confirmation

	err := h.UC.UploadAvatar(userAvatar.GetUserID(), userAvatar.GetFileName(), buf)
	_ = stream.SendAndClose(&empty.Empty{})
	return err
}

//функции вспомогательные для преобразования from/to models.User UserModel(GRPC)
func (h *UserServer) transfToModelsUser(in *UserModel) models.User {
	return models.User{
		ID:       in.GetID(),
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
		ID:       int64(in.ID),
		Password: in.Password,
		Name:     in.Name,
		Login:    in.Login,
		Image:    in.Image,
		Email:    in.Email,
	}
}
