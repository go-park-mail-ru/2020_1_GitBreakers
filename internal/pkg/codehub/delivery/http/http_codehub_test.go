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
	var limit int64 = 100
	var offset int64 = 0

	t.Run("StarredRepos ok", func(t *testing.T) {
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

	t.Run("StarredRepos wrong login", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			GetStarredRepos(gomock.AssignableToTypeOf(int64(0)), limit, offset).
			Return(repolist, nil).
			Times(0)
		UClientMock.EXPECT().
			GetByLogin(testUser.Login).
			Return(testUser, entityerrors.DoesNotExist()).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.StarredRepos, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"login": testUser.Login})

		apitest.New("StarredRepos wrong login").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + testUser.Login + "/stars").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})
	t.Run("StarredRepos err in getStarredRepos", func(t *testing.T) {
		gomock.InOrder(
			UClientMock.EXPECT().
				GetByLogin(testUser.Login).
				Return(testUser, nil).
				Times(1),
			UCCodeHubMock.EXPECT().
				GetStarredRepos(gomock.AssignableToTypeOf(int64(0)), limit, offset).
				Return(repolist, errors.New("some error")).
				Times(1),
		)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.StarredRepos, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"login": testUser.Login})

		apitest.New("StarredRepos err in getStarredRepos").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + testUser.Login + "/stars").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestCodeHubUserWithStar(t *testing.T) {
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

	const userCount int = 7
	userlist := make([]models.User, userCount)
	for i := range userlist {
		err := faker.FakeData(&userlist[i])
		require.Nil(t, err)
	}
	var limit int64 = 100
	var offset int64 = 0
	var RepoID int64 = 2
	RepoIDstr := "2"

	t.Run("UserWithStar ok", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			GetUserStaredList(RepoID, limit, offset).
			Return(userlist, nil).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.UserWithStar, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("UserWithStar ok").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + RepoIDstr + "/stars/users").
			Expect(t).
			Status(http.StatusOK).
			End()
	})

	t.Run("UserWithStar wrong mux vars", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			GetUserStaredList(RepoID, limit, offset).
			Return(userlist, nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.UserWithStar, false)

		apitest.New("UserWithStar wrong mux vars").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + RepoIDstr + "/stars/users").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})

	t.Run("UserWithStar repo not exsist", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			GetUserStaredList(RepoID, limit, offset).
			Return(userlist, entityerrors.DoesNotExist()).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.UserWithStar, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("UserWithStar repo not exsist").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + RepoIDstr + "/stars/users").
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("UserWithStar repo some err", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			GetUserStaredList(RepoID, limit, offset).
			Return(userlist, errors.New("some error")).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.UserWithStar, false)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("UserWithStar repo some err").
			Handler(middlewareMock).
			Method(http.MethodGet).
			URL("/func/repo/" + RepoIDstr + "/stars/users").
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestCodeHubNewIssue(t *testing.T) {
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

	const userCount int = 7
	userlist := make([]models.User, userCount)
	for i := range userlist {
		err := faker.FakeData(&userlist[i])
		require.Nil(t, err)
	}

	RepoIDstr := "2"
	someIssue := models.Issue{
		ID:        0,
		AuthorID:  12,
		RepoID:    2,
		Title:     "ffffff",
		Message:   "daw",
		Label:     "kek",
		IsClosed:  false,
		CreatedAt: time.Now().UTC(),
	}

	t.Run("NewIssue Unauthorized", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(someIssue).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, false)

		apitest.New("NewIssue Unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("NewIssue Bad repoID", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(someIssue).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": "kekID"})

		apitest.New("NewIssue Bad repoID").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("NewIssue Bad json body", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(someIssue).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("NewIssue Bad json body").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body("badJSON").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
	t.Run("NewIssue ok", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(gomock.AssignableToTypeOf(someIssue)).
			Return(nil).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("NewIssue ok").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, someIssue.Title, someIssue.Message)).
			Expect(t).
			Status(http.StatusCreated).
			End()
	})

	t.Run("NewIssue access denied", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(gomock.AssignableToTypeOf(someIssue)).
			Return(entityerrors.AccessDenied()).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("NewIssue access denied").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, someIssue.Title, someIssue.Message)).
			Expect(t).
			Status(http.StatusForbidden).
			End()
	})

	t.Run("NewIssue access denied", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(gomock.AssignableToTypeOf(someIssue)).
			Return(entityerrors.AlreadyExist()).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("NewIssue access denied").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, someIssue.Title, someIssue.Message)).
			Expect(t).
			Status(http.StatusConflict).
			End()
	})

	t.Run("NewIssue doesn't exsist repo", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(gomock.AssignableToTypeOf(someIssue)).
			Return(entityerrors.DoesNotExist()).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("NewIssue doesn't exsist repo").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, someIssue.Title, someIssue.Message)).
			Expect(t).
			Status(http.StatusNotFound).
			End()
	})

	t.Run("NewIssue some err in create", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			CreateIssue(gomock.AssignableToTypeOf(someIssue)).
			Return(errors.New("some error")).
			Times(1)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.NewIssue, true)
		middlewareMock = middleware.SetMuxVars(middlewareMock,
			map[string]string{"repoID": RepoIDstr})

		apitest.New("NewIssue some err in create").
			Handler(middlewareMock).
			Method(http.MethodPost).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body(fmt.Sprintf(`{"title": "%s", "message": "%s"}`, someIssue.Title, someIssue.Message)).
			Expect(t).
			Status(http.StatusInternalServerError).
			End()
	})
}
func TestCodeHubUpdateIssue(t *testing.T) {
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

	const userCount int = 7
	userlist := make([]models.User, userCount)
	for i := range userlist {
		err := faker.FakeData(&userlist[i])
		require.Nil(t, err)
	}

	RepoIDstr := "2"
	someIssue := models.Issue{
		ID:        0,
		AuthorID:  12,
		RepoID:    2,
		Title:     "ffffff",
		Message:   "daw",
		Label:     "kek",
		IsClosed:  false,
		CreatedAt: time.Now().UTC(),
	}

	t.Run("UpdIssue Unauthorized", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			UpdateIssue(someIssue).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.UpdateIssue, false)

		apitest.New("UpdIssue Unauthorized").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Expect(t).
			Status(http.StatusUnauthorized).
			End()
	})

	t.Run("UpdIssue bad json", func(t *testing.T) {
		UCCodeHubMock.EXPECT().
			UpdateIssue(someIssue).
			Return(nil).
			Times(0)

		middlewareMock := middleware.AuthMiddlewareMock(CodeHubHandlers.UpdateIssue, true)

		apitest.New("UpdIssue bad json").
			Handler(middlewareMock).
			Method(http.MethodPut).
			URL("/func/repo/" + RepoIDstr + "/issues").
			Body("kek json").
			Expect(t).
			Status(http.StatusBadRequest).
			End()
	})
}
