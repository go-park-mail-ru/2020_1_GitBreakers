package http

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	mock_clients "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/app/clients/mocks"
	mockCodehub "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
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

func TestCodeHubModifyStar(t *testing.T) {
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
	t.Run("Modify star bad repoid", func(t *testing.T) {

		repoID := "kek"

		UCCodeHubMock.EXPECT().
			ModifyStar(gomock.AssignableToTypeOf(star)).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.ModifyStar, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": repoID})

		apitest.New("Modify star bad repoid").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + repoID + "/stars").
			Body(fmt.Sprintf(`{"vote": true }`)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("Modify star bad json", func(t *testing.T) {

		repoID := strconv.Itoa(int(star.RepoID))

		UCCodeHubMock.EXPECT().
			ModifyStar(gomock.AssignableToTypeOf(star)).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.ModifyStar, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": repoID})

		apitest.New("Modify star bad json").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + repoID + "/stars").
			Body(fmt.Sprintf(`{vote": true }`)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("Modify star repo doesn't exsist", func(t *testing.T) {

		repoID := strconv.Itoa(int(star.RepoID))

		UCCodeHubMock.EXPECT().
			ModifyStar(gomock.AssignableToTypeOf(star)).
			Return(entityerrors.DoesNotExist()).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.ModifyStar, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": repoID})

		apitest.New("Modify star repo doesn't exsist").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + repoID + "/stars").
			Body(fmt.Sprintf(`{"vote": true }`)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})
	t.Run("Modify star repo some err", func(t *testing.T) {

		repoID := strconv.Itoa(int(star.RepoID))

		UCCodeHubMock.EXPECT().
			ModifyStar(gomock.AssignableToTypeOf(star)).
			Return(errors.New("some error")).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.ModifyStar, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": repoID})

		apitest.New("Modify star repo some err").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + repoID + "/stars").
			Body(fmt.Sprintf(`{"vote": true }`)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestCodeHubStarredRepos(t *testing.T) {
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

	const repoCount int = 7
	repolist := make([]gitmodels.Repository, repoCount)
	for i := range repolist {
		err := faker.FakeData(&repolist[i])
		require.Nil(t, err)
	}

	t.Run("StarredRepos ok", func(t *testing.T) {
		var limit int64 = 100
		var offset int64 = 0

		UCCodeHubMock.EXPECT().
			GetStarredRepos(gomock.AssignableToTypeOf(int64(0)), limit, offset).
			Return(repolist, nil).
			Times(1)
		UClientMock.EXPECT().
			GetByLogin(testUser.Login).
			Return(testUser, nil).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.StarredRepos, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"login": testUser.Login})

		apitest.New("StarredRepos ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + testUser.Login + "/stars").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
}
