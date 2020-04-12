package delivery

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	gitMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	userMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
)

var testRepo = gitmodels.Repository{
	ID:          0,
	OwnerID:     0,
	Name:        "djangoBlog",
	Description: "project to work with python",
	IsFork:      false,
	IsPublic:    true,
}
var gitHandlers GitDelivery

func TestGitDelivery_CreateRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	u := userMock.NewMockUCUser(ctrl)
	m := gitMock.NewMockUseCase(ctrl)

	gitHandlers.UserUC = u
	gitHandlers.UC = m
	gitHandlers.Logger = &newlogger

	t.Run("Create unauthorized", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().Create(gomock.Any(), testRepo).
				Return(nil).
				Times(0), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.CreateRepo, false)

		apitest.New("Create unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/repo").
			Body(fmt.Sprintf(`{"name": "%s", "description": "%s"}`, testRepo.Name, testRepo.Description)).
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("Create invalid json", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().Create(gomock.Any(), testRepo).
				Return(nil).
				Times(0), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.CreateRepo, true)

		apitest.New("Create invalid json").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/repo").
			Body(fmt.Sprintf(`{name": "%s", "description": "%s"}`, testRepo.Name, testRepo.Description)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Create invalid data", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().Create(gomock.Any(), testRepo).
				Return(nil).
				Times(0), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.CreateRepo, true)

		apitest.New("Create invalid data").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/repo").
			Body(fmt.Sprintf(`{"name": "%s", "description": "%s"}`, "", testRepo.Description)).
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("Create already exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().Create(gomock.Any(), gomock.Any()).
				Return(entityerrors.AlreadyExist()).
				Times(1), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.CreateRepo, true)

		apitest.New("Create already exsist").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/repo").
			Body(fmt.Sprintf(`{ "name": "%s", "description": "%s", "created_ad": "%s" }`,
				testRepo.Name, testRepo.Description, testRepo.CreatedAt)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})

	t.Run("Create some error", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().Create(gomock.Any(), gomock.Any()).
				Return(errors.New("some error")).
				Times(1), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.CreateRepo, true)

		apitest.New("Create already exsist").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/repo").
			Body(fmt.Sprintf(`{ "name": "%s", "description": "%s", "created_ad": "%s" }`,
				testRepo.Name, testRepo.Description, testRepo.CreatedAt)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
	t.Run("Create ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().Create(gomock.Any(), gomock.Any()).
				Return(nil).
				Times(1), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.CreateRepo, true)

		apitest.New("Create ok").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/repo").
			Body(fmt.Sprintf(`{ "name": "%s", "description": "%s", "created_ad": "%s" }`,
				testRepo.Name, testRepo.Description, testRepo.CreatedAt)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})
}
func TestUserHttp_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	u := userMock.NewMockUCUser(ctrl)
	m := gitMock.NewMockUseCase(ctrl)

	gitHandlers.UserUC = u
	gitHandlers.UC = m
	gitHandlers.Logger = &newlogger

	const repoCount int = 7
	repolist := make([]gitmodels.Repository, repoCount)
	for i := range repolist {
		err := faker.FakeData(&repolist[i])

		require.Nil(t, err)
	}

	t.Run("Get repolist ok", func(t *testing.T) {

		someLogin := "keksik"
		gomock.InOrder(
			m.EXPECT().GetRepoList(someLogin, gomock.Any()).
				Return(repolist, nil).
				Times(1), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepoList, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)

		apitest.New("Get repolist ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/repolist/keksik").
			Body("").
			Expect(t).
			Status(http.StatusOK).
			End()
	})
	t.Run("Get repolist access denied", func(t *testing.T) {

		someLogin := "keksik"
		gomock.InOrder(
			m.EXPECT().GetRepoList(someLogin, gomock.Any()).
				Return(repolist, entityerrors.AccessDenied()).
				Times(1), )

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepoList, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)

		apitest.New("Get repolist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/repolist/keksik").
			Body("").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})
}
