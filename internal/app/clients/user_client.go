package clients

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	usergrpc "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery/grpc"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type UserClient struct {
	conn   *grpc.ClientConn
	client usergrpc.UserGrpcClient
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

	copier.Copy(&userModelGRPC, &User)
	userModel := usergrpc.UserModel{
		Id:       int64(User.ID),
		Password: User.Password,
		Name:     User.Name,
		Login:    User.Login,
		Image:    User.Image,
		Email:    User.Email,
	}
	_, err := c.client.Create(context.Background(), &userModel)
	return err
}

func (c *UserClient) Update(userID int, newUserData models.User) error {
	userUPDModel := usergrpc.UserUpdateModel{
		UserID: int64(userID),
		UserData: &usergrpc.UserModel{
			Id:       0,
			Password: newUserData.Password,
			Name:     newUserData.Name,
			Login:    newUserData.Login,
			Image:    newUserData.Image,
			Email:    newUserData.Email,
		},
	}
	_, err := c.client.UpdateUser(context.Background(), &userUPDModel)
	return err
}
