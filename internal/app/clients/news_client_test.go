package clients

import (
	"context"
	"github.com/bxcodec/faker/v3"
	news "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/delivery/grpc"
	mockCodehub "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"testing"
	"time"
)

//мокаем usecase и тестим клиент->сервер->mockUseCase;
func TestNewsClient(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUCCodeHub := mockCodehub.NewMockUCCodeHubI(ctrl) //mockUCCodeHub для usecase
	news.NewNewsServer(server, mockUCCodeHub)            //кидаем этот мок в grpc server
	go func() {                                          //слушаем in-memory port
		if err := server.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	conn, _ := grpc.DialContext(context.Background(),
		"", grpc.WithContextDialer(bufDialer), grpc.WithInsecure())

	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	testUser := models.User{
		ID:        12,
		Password:  "123456789",
		Name:      "keks",
		Login:     "keksik",
		Image:     "default/image.jpg",
		Email:     "test@mail.ru",
		CreatedAt: time.Time{},
	}
	const newsCount int = 7
	newslist := make([]models.News, newsCount)
	for i := range newslist {
		err := faker.FakeData(&newslist[i].ID)
		newslist[i].Date = newslist[i].Date.UTC()
		require.Nil(t, err)
		err = faker.FakeData(&newslist[i].AuthorImage)
		require.Nil(t, err)
		err = faker.FakeData(&newslist[i].AuthorLogin)
		require.Nil(t, err)
		err = faker.FakeData(&newslist[i].Mess)
		require.Nil(t, err)
		err = faker.FakeData(&newslist[i].RepoID)
		require.Nil(t, err)

	}

	client := NewsClient{conn: conn, client: news.NewNewsClient(conn)}

	t.Run("Get OK", func(t *testing.T) {
		var someRepoID int64 = 5
		var limit int64 = 100
		var offset int64 = 0

		mockUCCodeHub.EXPECT().
			GetNews(someRepoID, testUser.ID, limit, offset).
			Return(newslist, nil).
			Times(1)

		newsSetFromDb, err := client.GetNews(someRepoID, testUser.ID, limit, offset)
		require.EqualValues(t, newslist, newsSetFromDb)
		require.Nil(t, err)
	})

	t.Run("New client ", func(t *testing.T) {
		client, err := NewNewsClient()
		require.Nil(t, err)
		require.NotNil(t, client)
	})

}
