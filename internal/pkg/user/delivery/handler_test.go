package delivery

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	sessMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/mocks"
	userMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"net/http"
	"os"
	"testing"
)

var testUser = models.User{
	ID:       34,
	Password: "52jkfgit389535dfe3",
	Name:     "somename",
	Login:    "dimaPetyaVasya",
	Image:    "default.png",
	Email:    "bezbab@mail.ru",
}
var userHandlers UserHttp

func TestCheckAuth(t *testing.T) {

	t.Run("Login-OK", func(t *testing.T) {

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := userMock.NewMockUCUser(ctrl)
		s := sessMock.NewMockSessDelivery(ctrl)

		userHandlers.UserUC = m
		newlogger := logger.NewTextFormatSimpleLogger(os.Stdout)
		userHandlers.Logger = &newlogger

		testInput := models.SignInForm{
			Login:    testUser.Login,
			Password: testUser.Password,
		}

		m.EXPECT().
			GetByLogin(testInput.Login).
			Return(testUser, nil)

		m.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil)

		s.EXPECT().
			Create(testUser.ID).
			Return(http.Cookie{
				Name:  "session_id",
				Value: "tj38r39i3r3j4953",
			}, nil)

		userHandlers.UserUC = m
		userHandlers.SessHttp = s

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false)

		apitest.New("Login-OK").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

}
