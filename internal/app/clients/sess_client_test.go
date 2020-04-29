package clients

import (
	"context"
	"errors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	session "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/delivery/grpc"
	sessMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"log"
	"testing"
	"time"
)

//мокаем usecase и тестим клиент->сервер->mockUseCase;
func TestSessClient(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mock := sessMock.NewMockUCSession(ctrl) //mock для usecase
	session.NewSessServer(server, mock)

	go func() { //слушаем in-memory port
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

	testSession := models.Session{
		ID:     "",
		UserID: 20,
	}

	client := SessClient{conn: conn, client: session.NewSessionClient(conn)}

	someErr := errors.New("some error")

	t.Run("CreateSess-OK", func(t *testing.T) {
		duration, err := time.ParseDuration("48h")
		require.NoError(t, err)

		someUUID := "123e4567-e89b-12d3-a456-426655440000"

		mock.EXPECT().
			Create(testSession, duration).
			Return(someUUID, nil).
			Times(1)

		sessUUID, err := client.CreateSess(testSession.UserID)

		require.NoError(t, err)
		require.Equal(t, sessUUID, someUUID)
	})

	t.Run("CreateSess-OK", func(t *testing.T) {
		duration, err := time.ParseDuration("48h")
		require.NoError(t, err)

		someUUID := "123e4567-e89b-12d3-a456-426655440000"

		mock.EXPECT().
			Create(testSession, duration).
			Return(someUUID, nil).
			Times(1)

		sessUUID, err := client.CreateSess(testSession.UserID)

		require.NoError(t, err)
		require.Equal(t, sessUUID, someUUID)
	})

	t.Run("CreateSess err", func(t *testing.T) {
		duration, err := time.ParseDuration("48h")
		require.NoError(t, err)

		someUUID := "123e4567-e89b-12d3-a456-426655440000"

		mock.EXPECT().
			Create(testSession, duration).
			Return(someUUID, someErr).
			Times(1)

		sessUUID, err := client.CreateSess(testSession.UserID)

		require.Error(t, err)
		errGrpc := status.Convert(err)

		require.Empty(t, sessUUID)
		require.EqualError(t, someErr, errGrpc.Message())
	})

	t.Run("DelSession full", func(t *testing.T) {

		someUUID := "123e4567-e89b-12d3-a456-426655440000"

		mock.EXPECT().
			Delete(someUUID).
			Return(someErr).
			Times(1)

		err := client.DelSess(someUUID)

		require.Error(t, err)
		errGrpc := status.Convert(err)

		require.EqualError(t, someErr, errGrpc.Message())
	})

	t.Run("GetSession ok", func(t *testing.T) {

		someUUID := "123e4567-e89b-12d3-a456-426655440000"
		testSession.ID = someUUID
		mock.EXPECT().
			GetByID(someUUID).
			Return(testSession, nil).
			Times(1)

		sessModel, err := client.GetSess(someUUID)

		require.NoError(t, err)

		require.Equal(t, sessModel, testSession)
	})

	t.Run("GetSession err", func(t *testing.T) {

		someUUID := "123e4567-e89b-12d3-a456-426655440000"
		testSession.ID = someUUID
		mock.EXPECT().
			GetByID(someUUID).
			Return(testSession, someErr).
			Times(1)

		sessModel, err := client.GetSess(someUUID)

		require.Error(t, err)

		require.Empty(t, sessModel)
	})
}
