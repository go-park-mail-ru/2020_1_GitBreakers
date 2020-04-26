package delivery

import (
	"errors"
	"fmt"
	"github.com/bxcodec/faker/v3"
	gitMock "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/middleware"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
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
				Times(0))

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
				Times(0))

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
				Times(0))

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
				Times(1))

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
				Times(1))

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
				Times(1))

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
func TestGitDelivery_GetRepoList(t *testing.T) {
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
				Times(1))

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
				Times(1))

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

	t.Run("Get repolist user empty", func(t *testing.T) {

		someUser := models.User{
			ID:       10,
			Password: "g4544tie48k34",
			Name:     "keksik",
			Login:    "vaper",
			Image:    "pudge.png",
			Email:    "dkfiksikrigb@mail.ru",
		}
		gomock.InOrder(
			u.EXPECT().GetByID(gomock.Any()).
				Return(someUser, nil).
				Times(1),
			m.EXPECT().GetRepoList(someUser.Login, gomock.Any()).
				Return(repolist, nil).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepoList, true)
		//middlewareMock = middleware.SetMuxVars(middlewareMock, "username", "")

		apitest.New("Get repolist user empty").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/repolist").
			Body("").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get repolist user not exsist", func(t *testing.T) {

		someUser := models.User{
			ID:       10,
			Password: "g4544tie48k34",
			Name:     "keksik",
			Login:    "vaper",
			Image:    "pudge.png",
			Email:    "dkfiksikrigb@mail.ru",
		}
		gomock.InOrder(
			u.EXPECT().GetByID(gomock.Any()).
				Return(someUser, errors.New("some error")).
				Times(1),
			m.EXPECT().GetRepoList(someUser.Login, gomock.Any()).
				Return(repolist, nil).
				Times(0))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepoList, true)

		apitest.New("Get repolist user not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/repolist").
			Body("").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
}
func TestGitDelivery_GetRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	u := userMock.NewMockUCUser(ctrl)
	m := gitMock.NewMockUseCase(ctrl)

	gitHandlers.UserUC = u
	gitHandlers.UC = m
	gitHandlers.Logger = &newlogger

	repo := gitmodels.Repository{}
	err := faker.FakeData(&repo)

	require.Nil(t, err)
	someLogin := "keksik"
	someRepo := "cheese"

	t.Run("Get repo ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().GetRepo(gomock.Eq(someLogin), gomock.Any(), gomock.Any()).
				Return(repo, nil).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepo, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get repo ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo).
			Body("").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get repo access denied", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().GetRepo(gomock.Eq(someLogin), gomock.Any(), gomock.Any()).
				Return(repo, entityerrors.AccessDenied()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepo, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get repo ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo).
			Body("").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Get repo not exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().GetRepo(gomock.Eq(someLogin), gomock.Any(), gomock.Any()).
				Return(repo, entityerrors.DoesNotExist()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepo, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get repo not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo).
			Body("").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Get repo some error", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().GetRepo(gomock.Eq(someLogin), gomock.Any(), gomock.Any()).
				Return(repo, errors.New("some error")).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetRepo, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get repo some error").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo).
			Body("").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestGitDelivery_GetBranchList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	u := userMock.NewMockUCUser(ctrl)
	m := gitMock.NewMockUseCase(ctrl)

	gitHandlers.UserUC = u
	gitHandlers.UC = m
	gitHandlers.Logger = &newlogger

	const branchCount int = 7
	branchlist := make([]gitmodels.Branch, branchCount)
	for i := range branchlist {
		err := faker.FakeData(&branchlist[i])

		require.Nil(t, err)
	}
	someLogin := "keksik"
	someRepo := "pepsertt"

	t.Run("Get branchlist ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetBranchList(nil, someLogin, gomock.Any()).
				Return(branchlist, nil).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetBranchList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get branchlist ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/branches").
			Body("").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get branchlist access denied", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetBranchList(nil, someLogin, gomock.Any()).
				Return(branchlist, entityerrors.AccessDenied()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetBranchList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get branchlist access denied").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/branches").
			Body("").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Get branchlist not exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetBranchList(nil, someLogin, gomock.Any()).
				Return(branchlist, entityerrors.DoesNotExist()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetBranchList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get branchlist not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/branches").
			Body("").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Get branchlist some error", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetBranchList(nil, someLogin, gomock.Any()).
				Return(branchlist, errors.New("some error")).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetBranchList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)

		apitest.New("Get branchlist some error").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/branches").
			Body("").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

func TestGitDelivery_GetCommitsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	u := userMock.NewMockUCUser(ctrl)
	m := gitMock.NewMockUseCase(ctrl)

	gitHandlers.UserUC = u
	gitHandlers.UC = m
	gitHandlers.Logger = &newlogger

	const commitsCount int = 7
	commitslist := make([]gitmodels.Commit, commitsCount)

	for i := range commitslist {
		err := faker.FakeData(&commitslist[i])

		require.Nil(t, err)
	}
	someLogin := "keksik"
	someRepo := "pepsertt"
	someBranch := "gdatwfw52g34"

	t.Run("Get commitslist ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByCommitHash(gomock.Any(), gomock.Any()).
				Return(commitslist, nil).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/commits/" + someBranch).
			Body("").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get commitslist forbidden", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByCommitHash(gomock.Any(), gomock.Any()).
				Return(commitslist, entityerrors.AccessDenied()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist forbidden").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/commits/" + someBranch).
			Body("").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Get commitslist not exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByCommitHash(gomock.Any(), gomock.Any()).
				Return(commitslist, entityerrors.DoesNotExist()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/commits/" + someBranch).
			Body("").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Get commitslist some err", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByCommitHash(gomock.Any(), gomock.Any()).
				Return(commitslist, errors.New("some error")).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsList, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist some err").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/commits/" + someBranch).
			Body("").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}

func TestGitDelivery_GetCommitsByBranchName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	newlogger := logger.NewTextFormatSimpleLogger(ioutil.Discard)
	u := userMock.NewMockUCUser(ctrl)
	m := gitMock.NewMockUseCase(ctrl)

	gitHandlers.UserUC = u
	gitHandlers.UC = m
	gitHandlers.Logger = &newlogger

	const commitsCount int = 7
	commitslist := make([]gitmodels.Commit, commitsCount)

	for i := range commitslist {
		err := faker.FakeData(&commitslist[i])

		require.Nil(t, err)
	}
	someLogin := "keksik"
	someRepo := "pepsertt"
	someBranch := "master"

	t.Run("Get commitslist ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByBranchName(gomock.Eq(someLogin), gomock.Any(),
					gomock.Any(), 0, 100, gomock.Any()).
				Return(commitslist, nil).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsByBranchName, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/" + someBranch + "/commits").
			Body("").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("Get commitslist forbidden", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByBranchName(gomock.Eq(someLogin), gomock.Any(),
					gomock.Any(), 0, 100, gomock.Any()).
				Return(commitslist, entityerrors.AccessDenied()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsByBranchName, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist forbidden").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/" + someBranch + "/commits").
			Body("").
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("Get commitslist not exsist", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByBranchName(gomock.Eq(someLogin), gomock.Any(),
					gomock.Any(), 0, 100, gomock.Any()).
				Return(commitslist, entityerrors.DoesNotExist()).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsByBranchName, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/" + someBranch + "/commits").
			Body("").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("Get commitslist some err", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				GetCommitsByBranchName(gomock.Eq(someLogin), gomock.Any(),
					gomock.Any(), 0, 100, gomock.Any()).
				Return(commitslist, errors.New("some error")).
				Times(1))

		middlewareMock := middleware.AuthMiddlewareMock(gitHandlers.GetCommitsByBranchName, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "username", someLogin)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "reponame", someRepo)
		middlewareMock = middleware.SetMuxVars(middlewareMock, "branchname", someBranch)

		apitest.New("Get commitslist some err").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/" + someLogin + "/" + someRepo + "/" + someBranch + "/commits").
			Body("").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
