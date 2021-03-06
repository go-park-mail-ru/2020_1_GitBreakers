package clients

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	usergrpc "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery/grpc"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client usergrpc.UserGrpcClient
}

func NewUserClient() (UserClient, error) {
	userClient := UserClient{}
	if err := userClient.Connect(); err != nil {
		return userClient, err
	}
	return userClient, nil
}

func (c *UserClient) Connect() error {
	//todo port unhardcode
	conn, err := grpc.Dial(":8082", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "grpc.Dial()")
	}
	c.conn = conn
	c.client = usergrpc.NewUserGrpcClient(c.conn)
	return nil
}
func (c *UserClient) Create(User models.User) error {
	userModelGRPC := usergrpc.UserModel{}

	if err := copier.Copy(&userModelGRPC, &User); err != nil {
		return err
	}

	_, err := c.client.Create(context.Background(), &userModelGRPC)

	stat := status.Convert(err)
	switch {
	case stat.Message() == entityerrors.AlreadyExist().Error():
		return entityerrors.AlreadyExist()
	default:
		return err
	}
}

func (c *UserClient) Update(userID int64, newUserData models.User) error {
	grpcUserModel := usergrpc.UserModel{}
	if err := copier.Copy(&grpcUserModel, &newUserData); err != nil {
		return err
	}
	userUPDModel := usergrpc.UserUpdateModel{
		UserID:   userID,
		UserData: &grpcUserModel,
	}
	_, err := c.client.UpdateUser(context.Background(), &userUPDModel)
	stat := status.Convert(err)
	switch {
	case stat.Message() == entityerrors.AlreadyExist().Error():
		return entityerrors.AlreadyExist()
	default:
		return err
	}
}
func (c *UserClient) GetByLogin(login string) (models.User, error) {
	loginGRPC := &usergrpc.LoginModel{Login: login}

	userGRPCModel, err := c.client.GetByLogin(context.Background(), loginGRPC)
	stat := status.Convert(err)
	switch {
	case stat.Message() == entityerrors.DoesNotExist().Error():
		return models.User{}, entityerrors.DoesNotExist()
	}

	userFromServer := models.User{}

	if err := copier.Copy(&userFromServer, userGRPCModel); err != nil {
		return models.User{}, err
	}
	return userFromServer, err
}

func (c *UserClient) GetByID(userID int64) (models.User, error) {
	idGRPC := &usergrpc.UserIDModel{UserID: userID}

	userGRPCModel, err := c.client.GetByID(context.Background(), idGRPC)
	stat := status.Convert(err)
	switch {
	case stat.Message() == entityerrors.DoesNotExist().Error():
		return models.User{}, entityerrors.DoesNotExist()
	}

	userFromServer := models.User{}

	if err := copier.Copy(&userFromServer, userGRPCModel); err != nil {
		return models.User{}, err
	}
	return userFromServer, err
}
func (c *UserClient) CheckPass(login string, pass string) (bool, error) {
	loginWithPassGRPC := &usergrpc.CheckPassModel{Login: login, Pass: pass}

	checkPassResp, err := c.client.CheckPass(context.Background(), loginWithPassGRPC)
	if checkPassResp != nil {
		return checkPassResp.GetIsCorrect(), err
	}
	stat := status.Convert(err)
	switch {
	case stat.Message() == entityerrors.DoesNotExist().Error():
		// если такого пользователя не существует
		return false, entityerrors.DoesNotExist()
	default:
		//в случае неуспешного запроса(ошибка клиента или сервера)
		return false, err
	}
}
func (c *UserClient) UploadAvatar(UserID int64, fileName string, fileData []byte, fileSize int64) (err error) {
	const ChunkSize int = 1 << 16 //64kb
	if int64(len(fileData)) != fileSize {
		return errors.New("can not assign real fileLen and received len")
	}
	stream, err := c.client.UploadAvatar(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		if _, closeErr := stream.CloseAndRecv(); closeErr != nil {
			if err != nil {
				err = errors.WithMessage(err, closeErr.Error())
			} else {
				err = closeErr
			}
		}
	}()
	//открываем поток и туда кидаем частями, удобно для передачи жирных картинок
	buf := make([]byte, ChunkSize)
	for i := 0; i < int(fileSize); i += ChunkSize {
		if i+ChunkSize < int(fileSize) {
			copy(buf, fileData[i:i+ChunkSize])
		} else {
			err = stream.Send(&usergrpc.UserAvatarModel{
				UserID:   UserID,
				FileName: fileName,
				Chunk:    fileData[i:],
			})
			if err != nil {
				return
			}
			break
		}

		err = stream.Send(&usergrpc.UserAvatarModel{
			UserID:   UserID,
			FileName: fileName,
			Chunk:    buf,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
