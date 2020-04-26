package usecase

import (
	"errors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestSessionUC_Create(t *testing.T) {

	t.Run("Create ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockSessRepo(ctrl)

		someUserID := 20
		someExpireTime := time.Duration(20) * time.Hour
		someRetSessionID := "base+f33g339ifk3933435"

		newSession := models.Session{
			ID:     "",
			UserID: someUserID,
		}
		m.EXPECT().
			Create(newSession, someExpireTime).
			Return(someRetSessionID, nil).Times(1)

		useCase := SessionUC{m}

		sessID, err := useCase.Create(newSession, someExpireTime)

		require.Equal(t, sessID, someRetSessionID)

		require.Nil(t, err)
	})

	t.Run("Create fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockSessRepo(ctrl)

		someUserID := 20
		someExpireTime := time.Duration(20) * time.Hour
		someRetSessionID := "base+f33g339ifk3933435"

		newSession := models.Session{
			ID:     "",
			UserID: someUserID,
		}

		someErr := errors.New("some error")

		m.EXPECT().
			Create(newSession, someExpireTime).
			Return(someRetSessionID, someErr).
			Times(1)

		useCase := SessionUC{m}

		sessID, err := useCase.Create(newSession, someExpireTime)

		require.Equal(t, sessID, someRetSessionID)

		require.Error(t, err)
	})
}
func TestSessionUC_Delete(t *testing.T) {

	t.Run("Delete ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockSessRepo(ctrl)

		someInputSessionID := "base+f33g339ifk3933435"

		m.EXPECT().
			DeleteByID(someInputSessionID).
			Return(nil).
			Times(1)

		useCase := SessionUC{m}

		err := useCase.Delete(someInputSessionID)

		require.Nil(t, err)
	})

	t.Run("Delete fail", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockSessRepo(ctrl)

		someInputSessionID := "base+f33g339ifk3933435"

		m.EXPECT().
			DeleteByID(someInputSessionID).
			Return(errors.New("some error")).
			Times(1)

		useCase := SessionUC{m}

		err := useCase.Delete(someInputSessionID)

		require.Error(t, err)
	})
}

func TestSessionUC_GetByID(t *testing.T) {

	t.Run("GetByID ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockSessRepo(ctrl)

		someInputSessionID := "base+f33g339ifk3933435"
		someUserID := 20

		someOldSession := models.Session{
			ID:     "",
			UserID: someUserID,
		}

		m.EXPECT().
			GetSessByID(someInputSessionID).
			Return(someOldSession, nil).
			Times(1)

		useCase := SessionUC{m}

		sessFromDB, err := useCase.GetByID(someInputSessionID)

		require.Equal(t, sessFromDB, someOldSession)

		require.NoError(t, err)
	})
}
