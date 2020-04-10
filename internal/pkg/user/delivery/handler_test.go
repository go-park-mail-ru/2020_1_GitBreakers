package delivery

import (
	"fmt"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	sessMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/mocks"
	userMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/steinfletcher/apitest"
	"io/ioutil"
	"net/http"
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

func TestUserHttp_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := userMock.NewMockUCUser(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)

	userHandlers.UserUC = m
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	userHandlers.Logger = &newlogger

	testInput := models.SignInForm{
		Login:    testUser.Login,
		Password: testUser.Password,
	}
	userHandlers.UserUC = m
	userHandlers.SessHttp = s

	t.Run("Login-OK", func(t *testing.T) {
		m.EXPECT().
			GetByLogin(testInput.Login).
			Return(testUser, nil).Times(1)
		m.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(1)

		s.EXPECT().
			Create(testUser.ID).
			Return(http.Cookie{
				Name:  "session_id",
				Value: "tj38r39i3r3j4953",
			}, nil).Times(1)

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

	t.Run("Login already auth", func(t *testing.T) {
		m.EXPECT().
			GetByLogin(testInput.Login).
			Return(testUser, nil).Times(0)
		m.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(0)

		s.EXPECT().
			Create(testUser.ID).
			Return(http.Cookie{
				Name:  "session_id",
				Value: "tj38r39i3r3j4953",
			}, nil).Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, true)

		apitest.New("Login already auth").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusNotAcceptable).
			End()
	})

	t.Run("User not exsist", func(t *testing.T) {
		m.EXPECT().
			GetByLogin(testInput.Login).
			Return(models.User{}, entityerrors.DoesNotExist()).Times(1)

		m.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(0)

		s.EXPECT().
			Create(testUser.ID).
			Return(http.Cookie{
				Name:  "session_id",
				Value: "tj38r39i3r3j4953",
			}, nil).Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false)

		apitest.New("User not exsist").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
	t.Run("Some error in UseCase", func(t *testing.T) {
		m.EXPECT().
			GetByLogin(testInput.Login).
			Return(models.User{}, errors.New("some error")).Times(1)

		m.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(0)


		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false)

		apitest.New("Some error in UseCase").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Error in session", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByLogin(testInput.Login).
				Return(testUser, nil).Times(1),
			m.EXPECT().
				CheckPass(testInput.Login, testInput.Password).
				Return(true, nil).Times(1),
			s.EXPECT().
				Create(testUser.ID).
				Return(http.Cookie{
					Name:  "session_id",
					Value: "tj38r39i3r3j4953",
				}, errors.New("some error")).Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false)

		apitest.New("Error in session").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Invalid json", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByLogin(testInput.Login).
				Return(testUser, nil).Times(0))

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false)

		apitest.New("Invalid json").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login: "%s", "password": "%s"}`, testUser.Login, testUser.Password)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Json ok, data invalid", func(t *testing.T) {
		invalidPassword := "45"
		invalidLogin := "kek"
		gomock.InOrder(
			m.EXPECT().
				GetByLogin(invalidLogin).
				Return(testUser, nil).Times(0))

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Login, false)

		apitest.New("Invalid json").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/login").
			Body(fmt.Sprintf(`{"login": "%s", "password": "%s"}`, invalidLogin, invalidPassword)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
