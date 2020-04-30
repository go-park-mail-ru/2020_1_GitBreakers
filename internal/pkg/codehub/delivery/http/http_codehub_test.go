package http

import (
	"fmt"
	mock_clients "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/mocks"
	mockCodehub "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"io/ioutil"
	"net/http"
	"strconv"
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
var CodeHubHandlers HttpCodehub

func TestUserHttp_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	UCCodeHubMock := mockCodehub.NewMockUCCodeHubI(ctrl)
	UClientMock := mock_clients.NewMockUserClientI(ctrl)
	NewsClientMock := mock_clients.NewMockNewsClientI(ctrl)
	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)

	CodeHubHandlers.Logger = &newlogger
	CodeHubHandlers.NewsClient = NewsClientMock
	CodeHubHandlers.UserClient = UClientMock
	CodeHubHandlers.CodeHubUC = UCCodeHubMock
	star := models.Star{
		AuthorID: 12,
		RepoID:   12,
		Vote:     true,
	}

	t.Run("Modify star ok", func(t *testing.T) {

		repoID := strconv.Itoa(int(star.RepoID))

		UCCodeHubMock.EXPECT().
			ModifyStar(gomock.AssignableToTypeOf(star)).
			Return(nil).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.ModifyStar, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": repoID})

		apitest.New("Modify star ok").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + repoID + "/stars").
			Body(fmt.Sprintf(`{"vote": true }`)).
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Modify star unauthorized", func(t *testing.T) {

		repoID := strconv.Itoa(int(star.RepoID))

		UCCodeHubMock.EXPECT().
			ModifyStar(gomock.AssignableToTypeOf(star)).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.ModifyStar, false)

		apitest.New("Modify star unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + repoID + "/stars").
			Body(fmt.Sprintf(`{"vote": true }`)).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})
}
