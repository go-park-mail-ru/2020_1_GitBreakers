package http

import (
	"fmt"
	"github.com/bxcodec/faker/v3"
	mock_clients "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	sessMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/session/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"io/ioutil"
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
var userHandlers UserHttp

func TestUserHttp_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	u := sessMock.NewMockSessDelivery(ctrl)
	s := mock_clients.NewMockUserClientI(ctrl)

	userHandlers.UClient = s
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)
	userHandlers.Logger = &newlogger

	testInput := models.SignInForm{
		Login:    testUser.Login,
		Password: testUser.Password,
	}
	userHandlers.SessHttp = u

	t.Run("Login-OK", func(t *testing.T) {
		s.EXPECT().
			GetByLogin(testInput.Login).
			Return(testUser, nil).Times(1)
		s.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(1)

		u.EXPECT().
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
		s.EXPECT().
			GetByLogin(testInput.Login).
			Return(testUser, nil).Times(0)
		s.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(0)

		u.EXPECT().
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
		s.EXPECT().
			GetByLogin(testInput.Login).
			Return(models.User{}, entityerrors.DoesNotExist()).Times(1)

		s.EXPECT().
			CheckPass(testInput.Login, testInput.Password).
			Return(true, nil).Times(0)

		u.EXPECT().
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
		s.EXPECT().
			GetByLogin(testInput.Login).
			Return(models.User{}, errors.New("some error")).Times(1)

		s.EXPECT().
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
			s.EXPECT().
				GetByLogin(testInput.Login).
				Return(testUser, nil).Times(1),
			s.EXPECT().
				CheckPass(testInput.Login, testInput.Password).
				Return(true, nil).Times(1),
			u.EXPECT().
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
			s.EXPECT().
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
			s.EXPECT().
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

func TestUserHttp_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_clients.NewMockUserClientI(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)

	userHandlers.UClient = m
	userHandlers.SessHttp = s
	userHandlers.Logger = &newlogger

	testInput := models.User{}
	err := faker.FakeData(&testInput)

	require.Nil(t, err)

	var testUserEmpty = models.User{
		Password: "52jkfgit389535dfe3",
		Name:     "",
		Login:    "dimaPetyaVasya",
		Image:    "",
		Email:    "bezbab@mail.ru",
	}

	t.Run("Signup already auth", func(t *testing.T) {
		m.EXPECT().
			Create(testInput).
			Return(nil).Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, true)

		apitest.New("Signup already auth").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUser.Email, testUser.Password, testUser.Login)).
			Expect(t).
			Status(http.StatusNotAcceptable).
			End()
	})
	t.Run("Signup invalid json", func(t *testing.T) {
		m.EXPECT().
			Create(testInput).
			Return(nil).Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup invalid json").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "emal": "%s", "password": "%s", "login": "%s" }`,
				testUser.Email, testUser.Password, testUser.Login)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Signup invalid data", func(t *testing.T) {
		m.EXPECT().
			Create(testInput).
			Return(nil).Times(0)
		bademail := "@bademail@"

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup invalid data").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				bademail, testUser.Password, testUser.Login)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Signup good create", func(t *testing.T) {
		sessionCookie := http.Cookie{
			Name:     "session-id",
			Value:    "25425fg3f3535",
			Path:     "/",
			Domain:   "89.208.198.186",
			Expires:  time.Now(),
			MaxAge:   0,
			Secure:   false,
			HttpOnly: false,
		}
		gomock.InOrder(
			m.EXPECT().
				Create(testUserEmpty).
				Return(nil).Times(1),
			m.EXPECT().
				GetByLogin(testUser.Login).
				Return(testUser, nil).Times(1),
			s.EXPECT().
				Create(gomock.Any()).
				Return(sessionCookie, nil).Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup good create").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUser.Email, testUser.Password, testUser.Login)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("Signup invalid json", func(t *testing.T) {

		gomock.InOrder(
			m.EXPECT().
				Create(testUser).
				Return(nil).Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup invalid json").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email: "%s", "password": "%s", "login": "%s" }`,
				testUser.Email, testUser.Password, testUser.Login)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Signup already exsist", func(t *testing.T) {

		gomock.InOrder(
			m.EXPECT().
				Create(testUserEmpty).
				Return(entityerrors.AlreadyExist()).Times(1),
			m.EXPECT().
				GetByLogin(testUser.Login).
				Return(testUser, nil).Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup already exsist").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})
	t.Run("Signup error in create", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				Create(testUserEmpty).
				Return(errors.New("some error")).
				Times(1),
			m.EXPECT().
				GetByLogin(testUser.Login).
				Return(testUser, nil).Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup error in create").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("Signup some err in getByLogin func", func(t *testing.T) {

		gomock.InOrder(
			m.EXPECT().
				Create(testUserEmpty).
				Return(nil).
				Times(1),
			m.EXPECT().
				GetByLogin(testUser.Login).
				Return(testUser, errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup some err in getByLogin func").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("Signup some err in Sess Create func", func(t *testing.T) {

		sessionCookie := http.Cookie{
			Name:     "session-id",
			Value:    "25425fg3f3535",
			Path:     "/",
			Domain:   "89.208.198.186",
			Expires:  time.Now(),
			MaxAge:   0,
			Secure:   false,
			HttpOnly: false,
		}

		gomock.InOrder(
			m.EXPECT().
				Create(testUserEmpty).
				Return(nil).
				Times(1),
			m.EXPECT().
				GetByLogin(testUser.Login).
				Return(testUser, nil).
				Times(1),
			s.EXPECT().
				Create(testUser.ID).
				Return(sessionCookie, errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Create, false)

		apitest.New("Signup some err in Sess Create func").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/signup").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

}
func TestUserHttp_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_clients.NewMockUserClientI(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)

	userHandlers.UClient = m
	userHandlers.SessHttp = s
	userHandlers.Logger = &newlogger

	testInput := models.User{}
	err := faker.FakeData(&testInput)

	require.Nil(t, err)

	var testUserEmpty = models.User{
		Password: "52jkfgit389535dfe3",
		Name:     "",
		Login:    "dimaPetyaVasya",
		Image:    "",
		Email:    "bezbab@mail.ru",
	}

	t.Run("Update unauthorized", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				Update(gomock.Eq(10), testUser).
				Return(nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, false)

		apitest.New("Update unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/profile").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Update ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				Update(gomock.AssignableToTypeOf(int64(0)), testUserEmpty).
				Return(nil).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true)

		apitest.New("Update ok").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/profile").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Update invalid json", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				Update(gomock.AssignableToTypeOf(int64(0)), testUserEmpty).
				Return(nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true)

		apitest.New("Update invalid json").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/profile").
			Body(fmt.Sprintf(`{ email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("Update already exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				Update(gomock.AssignableToTypeOf(int64(0)), testUserEmpty).
				Return(entityerrors.AlreadyExist()).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true)

		apitest.New("Update already exsist").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/profile").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})
	t.Run("Update some err in Update func", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				Update(gomock.AssignableToTypeOf(int64(0)), testUserEmpty).
				Return(errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Update, true)

		apitest.New("Update some err in Update func").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/profile").
			Body(fmt.Sprintf(`{ "email": "%s", "password": "%s", "login": "%s" }`,
				testUserEmpty.Email, testUserEmpty.Password, testUserEmpty.Login)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestUserHttp_Logout(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_clients.NewMockUserClientI(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)

	userHandlers.UClient = m
	userHandlers.SessHttp = s
	userHandlers.Logger = &newlogger

	testInput := models.User{}
	err := faker.FakeData(&testInput)

	require.Nil(t, err)

	t.Run("Logout unauthorized", func(t *testing.T) {
		gomock.InOrder(
			s.EXPECT().
				Delete(gomock.AssignableToTypeOf("")).
				Return(nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, false)

		apitest.New("Logout unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/logout").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Logout session not exsist", func(t *testing.T) {
		gomock.InOrder(
			s.EXPECT().
				Delete(gomock.AssignableToTypeOf("")).
				Return(nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, true)

		apitest.New("Logout session not exsist").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/logout").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
	t.Run("Logout ok", func(t *testing.T) {
		gomock.InOrder(
			s.EXPECT().
				Delete("4ri39r3u438r3u438fhj3f384jf3438").
				Return(nil).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, true)

		apitest.New("Logout unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/logout").
			Cookie("session_id", "4ri39r3u438r3u438fhj3f384jf3438").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Logout session not delete", func(t *testing.T) {
		gomock.InOrder(
			s.EXPECT().
				Delete("4ri39r3u438r3u438fhj3f384jf3438").
				Return(errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.Logout, true)

		apitest.New("Logout unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/logout").
			Cookie("session_id", "4ri39r3u438r3u438fhj3f384jf3438").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestUserHttp_GetInfo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_clients.NewMockUserClientI(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)

	userHandlers.UClient = m
	userHandlers.SessHttp = s
	userHandlers.Logger = &newlogger

	testInput := models.User{}
	err := faker.FakeData(&testInput)

	require.Nil(t, err)

	t.Run("Getinto unauthorized", func(t *testing.T) {
		gomock.InOrder(
			s.EXPECT().
				Delete(gomock.AssignableToTypeOf("")).
				Return(nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfo, false)

		apitest.New("Logout unauthorized").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/logout").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Getinto ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByID(gomock.Any()).
				Return(testInput, nil).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfo, true)

		apitest.New("Getinto ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/whoami").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Getinto entity does not exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByID(gomock.AssignableToTypeOf(int64(0))).
				Return(testInput, entityerrors.DoesNotExist()).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfo, true)

		apitest.New("Getinto entity does not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/whoami").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Getinto some error in DB", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByID(gomock.AssignableToTypeOf(int64(0))).
				Return(testInput, errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfo, true)

		apitest.New("Getinto entity does not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/whoami").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})

	t.Run("Getinto some error in DB", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByID(gomock.AssignableToTypeOf(int64(0))).
				Return(testInput, errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfo, true)

		apitest.New("Getinto some error in DB").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/whoami").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestUserHttp_UploadAvatar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_clients.NewMockUserClientI(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)

	userHandlers.UClient = m
	userHandlers.SessHttp = s
	userHandlers.Logger = &newlogger

	testInput := models.User{}
	err := faker.FakeData(&testInput)

	require.Nil(t, err)

	t.Run("Upload avatar empty image", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByID(gomock.AssignableToTypeOf(int64(0))).
				Return(testUser, nil).
				Times(0),
			m.EXPECT().
				UploadAvatar(gomock.AssignableToTypeOf(int64(0)),
					gomock.AssignableToTypeOf(""),
					gomock.AssignableToTypeOf([]byte("")),
					gomock.AssignableToTypeOf(int64(0))).
				Return(nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.UploadAvatar, true)

		apitest.New("Upload avatar unauthorized").
			Debug().
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("user/avatar").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Upload avatar unauthorized", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetByID(gomock.Any()).
				Return(testInput, nil).
				Times(0),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.UploadAvatar, false)

		apitest.New("Upload avatar unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/avatar").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
}
func TestUserHttp_GetInfoByLogin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_clients.NewMockUserClientI(ctrl)
	s := sessMock.NewMockSessDelivery(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard, 1)

	userHandlers.UClient = m
	userHandlers.SessHttp = s
	userHandlers.Logger = &newlogger

	testInput := models.User{}
	err := faker.FakeData(&testInput)

	require.Nil(t, err)

	t.Run("Get info by login", func(t *testing.T) {
		someLogin := "keksik"
		err := faker.FakeData(&someLogin)

		require.Nil(t, err)

		gomock.InOrder(
			m.EXPECT().
				GetByLogin(someLogin).
				Return(testUser, nil).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfoByLogin, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, map[string]string{"login": someLogin})

		apitest.New("Get info by login").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/profile/" + someLogin).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get info login doesn't exsist", func(t *testing.T) {
		someLogin := "keksik"
		err := faker.FakeData(&someLogin)

		require.Nil(t, err)

		gomock.InOrder(
			m.EXPECT().
				GetByLogin(someLogin).
				Return(models.User{}, entityerrors.DoesNotExist()).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfoByLogin, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, map[string]string{"login": someLogin})

		apitest.New("Get info login doesn't exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/profile/" + someLogin).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
	t.Run("Get info login some error", func(t *testing.T) {
		someLogin := "keksik"
		err := faker.FakeData(&someLogin)

		require.Nil(t, err)

		gomock.InOrder(
			m.EXPECT().
				GetByLogin(someLogin).
				Return(models.User{}, errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(userHandlers.GetInfoByLogin, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, map[string]string{"login": someLogin})

		apitest.New("Get info login doesn't exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/profile/" + someLogin).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
