package http

import (
	"errors"
	mockClients "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

var testUser = models.User{
	ID:       34,
	Password: "52jkfgit389535dfe3",
	Name:     "somename",
	Login:    "dimaPetyaVasya",
	Image:    "default.png",
	Email:    "bezbab@mail.ru",
}
var sessHandler SessionHttp

func TestSessionHttp_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mockClients.NewMockSessClientI(ctrl)

	sessHandler.Client = s

	t.Run("Delete OK", func(t *testing.T) {
		someSessID := "somefkw3942"

		s.EXPECT().
			DelSess(someSessID).
			Return(nil).
			Times(1)

		err := sessHandler.Delete(someSessID)

		require.NoError(t, err)
	})
}

func TestSessionHttp_GetBySessID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mockClients.NewMockSessClientI(ctrl)

	sessHandler.Client = s

	t.Run("Get OK", func(t *testing.T) {
		someSessID := "somefkw3942"
		var someUserID int64 = 20

		someSession := models.Session{
			ID:     someSessID,
			UserID: someUserID,
		}
		s.EXPECT().
			GetSess(someSessID).
			Return(someSession, nil).
			Times(1)

		sessFromDB, err := sessHandler.GetBySessID(someSessID)

		require.NoError(t, err)
		require.Equal(t, sessFromDB, someSession)
	})
}

func TestSessionHttp_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := mockClients.NewMockSessClientI(ctrl)

	sessHandler.Client = s

	t.Run("Create OK", func(t *testing.T) {
		someSessID := "somefkw3942"
		var someUserID int64 = 20

		emptySession := models.Session{
			ID:     "",
			UserID: someUserID,
		}

		fullSession := models.Session{
			ID:     someSessID,
			UserID: someUserID,
		}

		someCookie := http.Cookie{
			Name:     "session_id",
			Value:    fullSession.ID,
			HttpOnly: true,
			Expires:  time.Now().Add(sessHandler.ExpireTime),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   false,
		}

		s.EXPECT().
			CreateSess(emptySession.UserID).
			Return(fullSession.ID, nil).
			Times(1)

		cookieFromDB, err := sessHandler.Create(someUserID)

		require.NoError(t, err)
		require.Equal(t, someCookie.Value, cookieFromDB.Value)
	})
	t.Run("Create error", func(t *testing.T) {
		someSessID := "somefkw3942"
		var someUserID int64 = 20

		emptySession := models.Session{
			ID:     "",
			UserID: someUserID,
		}

		fullSession := models.Session{
			ID:     someSessID,
			UserID: someUserID,
		}

		someCookie := http.Cookie{
			Name:     "session_id",
			Value:    fullSession.ID,
			HttpOnly: true,
			Expires:  time.Now().Add(sessHandler.ExpireTime),
			Path:     "/",
			SameSite: http.SameSiteNoneMode,
			Secure:   false,
		}

		s.EXPECT().
			CreateSess(emptySession.UserID).
			Return(fullSession.ID, errors.New("some error")).
			Times(1)

		cookieFromDB, err := sessHandler.Create(someUserID)

		require.Error(t, err)
		require.NotEqual(t, someCookie.Value, cookieFromDB.Value)
	})
}
