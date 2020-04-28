package clients

import (
	"context"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	usergrpc "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/delivery/grpc"
	userMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
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

func TestUserCreate(t *testing.T) {
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

}
