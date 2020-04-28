package clients

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	usergrpc "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery/grpc"
	userMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"net"
	"testing"
	"time"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

//мокаем usecase и тестим клиент->сервер->mockUseCase;
func TestUserClient(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := userMock.NewMockUCUser(ctrl) //mock для usecase
	usergrpc.NewUserServer(server, mock) //кидаем этот мок в grpc server
	go func() {                          //слушаем in-memory port
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	conn, _ := grpc.DialContext(context.Background(),
		"", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())

	defer conn.Close()

	testUser := models.User{
		ID:        12,
		Password:  "123456789",
		Name:      "keks",
		Login:     "keksik",
		Image:     "default/image.jpg",
		Email:     "test@mail.ru",
		CreatedAt: time.Time{},
	}

	client := UserClient{conn: conn, client: usergrpc.NewUserGrpcClient(conn)}

	t.Run("Create-OK", func(t *testing.T) {
		mock.EXPECT().Create(testUser).Return(nil).Times(1)
		err := client.Create(testUser)
		require.Nil(t, err)
	})

	t.Run("Create some err", func(t *testing.T) {
		someErr := errors.New("some error")
		mock.EXPECT().Create(testUser).Return(someErr).Times(1)
		err := client.Create(testUser)

		require.NotEqual(t, err, someErr)

		grpcErr := status.Convert(err)

		require.Equal(t, err, grpcErr.Err())
	})

	t.Run("Get by login ok", func(t *testing.T) {
		mock.EXPECT().GetByLogin(testUser.Login).Return(testUser, nil).Times(1)
		userFromDB, err := client.GetByLogin(testUser.Login)

		require.Equal(t, userFromDB, testUser)

		grpcErr := status.Convert(err)

		require.Equal(t, err, grpcErr.Err())
	})

	t.Run("GetByLogin err", func(t *testing.T) {
		mock.EXPECT().GetByLogin(testUser.Login).Return(testUser, entityerrors.DoesNotExist()).Times(1)
		userFromDB, err := client.GetByLogin(testUser.Login)

		require.Equal(t, userFromDB, models.User{})

		grpcErr := status.Convert(err)

		require.EqualError(t, entityerrors.DoesNotExist(), grpcErr.Message())
	})

	t.Run("GetByID ok", func(t *testing.T) {
		mock.EXPECT().GetByID(testUser.ID).Return(testUser, nil).Times(1)
		userFromDB, err := client.GetByID(testUser.ID)

		require.Equal(t, userFromDB, testUser)

		grpcErr := status.Convert(err)

		require.Nil(t, grpcErr)
	})

	t.Run("GetByID err", func(t *testing.T) {
		mock.EXPECT().GetByID(testUser.ID).Return(testUser, entityerrors.DoesNotExist()).Times(1)
		userFromDB, err := client.GetByID(testUser.ID)

		require.Equal(t, userFromDB, models.User{})

		grpcErr := status.Convert(err)

		require.EqualError(t, entityerrors.DoesNotExist(), grpcErr.Message())

	})
	t.Run("CheckPass ok", func(t *testing.T) {
		mock.EXPECT().CheckPass(testUser.Login, testUser.Password).Return(true, nil).Times(1)
		is_correct, err := client.CheckPass(testUser.Login, testUser.Password)

		require.Equal(t, is_correct, true)

		require.Nil(t, err)

	})

}
