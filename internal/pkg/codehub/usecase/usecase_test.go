package usecase

import (
	"errors"
	"github.com/bxcodec/faker/v3"
	mockCodehub "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/mocks"
	mocksGit "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	perm "github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/permission_types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var someUserOwner = models.User{
	ID:       12,
	Password: "sjfsfser242df",
	Name:     "Kekkers",
	Login:    "alahahbar",
	Image:    "/static/image/avatar/kek.png",
	Email:    "putinka@kremlin.ru",
}
var someUserNoAccess = models.User{
	ID:       10,
	Password: "sjfsfser242df",
	Name:     "Kekkers",
	Login:    "hehehmda",
	Image:    "/static/image/avatar/kek.png",
	Email:    "putin@kremlin.ru",
}

var someRepo = gitmodels.Repository{
	ID:          45352,
	OwnerID:     12,
	Name:        "faffafsaf",
	Description: "fasfafafaf",
	IsFork:      false,
	CreatedAt: time.Date(
		2015, 11, 17, 20, 34, 58, 651387237, time.UTC),
	IsPublic: false,
	Stars:    2,
}

func TestUCCodeHubStar(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mockCodehub.NewMockRepoStarI(ctrl)

	useCase := UCCodeHub{
		RepoStar: m,
	}

	t.Run("Modify star +", func(t *testing.T) {
		m.EXPECT().AddStar(someUserOwner.ID, someRepo.ID).Return(nil)

		err := useCase.ModifyStar(models.Star{
			AuthorID: someUserOwner.ID,
			RepoID:   someRepo.ID,
			Vote:     true,
		})

		require.NoError(t, err)
	})

	t.Run("Modify star -", func(t *testing.T) {
		m.EXPECT().DelStar(someUserOwner.ID, someRepo.ID).Return(nil)

		err := useCase.ModifyStar(models.Star{
			AuthorID: someUserOwner.ID,
			RepoID:   someRepo.ID,
			Vote:     false,
		})

		require.NoError(t, err)
	})

	t.Run("Get starred repos", func(t *testing.T) {
		const repoCount int = 7
		repolist := make([]gitmodels.Repository, repoCount)
		for i := range repolist {
			err := faker.FakeData(&repolist[i])
			require.Nil(t, err)
		}
		limit := 10
		offset := 2
		m.EXPECT().GetStarredRepos(someUserOwner.ID, int64(limit), int64(offset)).Return(repolist, nil)

		reposetFromDb, err := useCase.GetStarredRepos(someUserOwner.ID, int64(limit), int64(offset))

		require.EqualValues(t, reposetFromDb, repolist) //тк репосет и []models.repository разные типы
		require.NoError(t, err)

	})

	t.Run("Get user starred list", func(t *testing.T) {
		const userCount int = 10
		userlist := make([]models.User, userCount)
		for i := range userlist {
			err := faker.FakeData(&userlist[i])
			require.Nil(t, err)
		}
		limit := 10
		offset := 2
		m.EXPECT().
			GetUserStaredList(someUserOwner.ID, int64(limit), int64(offset)).
			Return(userlist, nil).
			Times(1)

		reposetFromDb, err := useCase.GetUserStaredList(someUserOwner.ID, int64(limit), int64(offset))

		require.EqualValues(t, reposetFromDb, userlist) //тк репосет и []models.repository разные типы
		require.NoError(t, err)
	})

}
func TestUCCodeHubNews(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	repoNews := mockCodehub.NewMockRepoNewsI(ctrl)
	userRepo := mocks.NewMockRepoUser(ctrl)
	gitRepo := mocksGit.NewMockGitRepoI(ctrl)
	useCase := UCCodeHub{
		RepoNews: repoNews,
		UserRepo: userRepo,
		GitRepo:  gitRepo,
	}
	const newsCount int = 7
	newslist := make([]models.News, newsCount)
	for i := range newslist {
		err := faker.FakeData(&newslist[i])
		require.Nil(t, err)
	}
	t.Run("Get news", func(t *testing.T) {
		var limit int64 = 100
		var offset int64 = 5

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(&someUserOwner.ID, someUserOwner.Login, someRepo.Name).
				Return(true, nil).
				Times(1),
			repoNews.EXPECT().
				GetNews(someRepo.ID, limit, offset).
				Return(newslist, nil).
				Times(1),
		)

		newslistFromDB, err := useCase.GetNews(someRepo.ID, someUserOwner.ID, limit, offset)

		require.EqualValues(t, newslistFromDB, newslist)

		require.NoError(t, err)
	})
	t.Run("Get news invalid repoID", func(t *testing.T) {
		var limit int64 = 100
		var offset int64 = 5

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, entityerrors.DoesNotExist()).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(0),
		)

		newslistFromDB, err := useCase.GetNews(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Empty(t, newslistFromDB)

		require.Error(t, err)
	})
	t.Run("Get news", func(t *testing.T) {
		var limit int64 = 100
		var offset int64 = 5

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, entityerrors.DoesNotExist()).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(&someUserOwner.ID, someUserOwner.Login, someRepo.Name).
				Return(true, nil).
				Times(0),
		)

		newslistFromDB, err := useCase.GetNews(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Empty(t, newslistFromDB)

		require.Error(t, err)
	})
	t.Run("Get news err in read access", func(t *testing.T) {
		var limit int64 = 100
		var offset int64 = 5

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(&someUserOwner.ID, someUserOwner.Login, someRepo.Name).
				Return(false, errors.New("some error")).
				Times(1),
			repoNews.EXPECT().
				GetNews(someRepo.ID, limit, offset).
				Return(newslist, nil).
				Times(0),
		)

		newslistFromDB, err := useCase.GetNews(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Empty(t, newslistFromDB)

		require.Error(t, err)
	})
	t.Run("Get news error with access", func(t *testing.T) {
		var limit int64 = 100
		var offset int64 = 5

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(&someUserOwner.ID, someUserOwner.Login, someRepo.Name).
				Return(false, errors.New("some error")).
				Times(1),
		)

		newslistFromDB, err := useCase.GetNews(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Empty(t, newslistFromDB)

		require.Error(t, err)

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(&someUserOwner.ID, someUserOwner.Login, someRepo.Name).
				Return(false, nil).
				Times(1),
		)

		newslistFromDB, err = useCase.GetNews(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Empty(t, newslistFromDB)

		require.Equal(t, err, entityerrors.AccessDenied())
	})
}
func TestUCCodeHubIssues(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repoIssues := mockCodehub.NewMockRepoIssueI(ctrl)
	userRepo := mocks.NewMockRepoUser(ctrl)
	gitRepo := mocksGit.NewMockGitRepoI(ctrl)

	useCase := UCCodeHub{
		RepoIssue: repoIssues,
		GitRepo:   gitRepo,
		UserRepo:  userRepo,
	}
	someIssue := models.Issue{
		ID:       20,
		AuthorID: 12,
		RepoID:   45352,
		Title:    "how to increase money",
		Message:  "no money",
		Label:    "fixed",
		IsClosed: false,
		CreatedAt: time.Date(
			2017, 11, 17, 20, 34, 58, 651387237, time.UTC),
	}
	const issuesCount int = 10
	issueslist := make([]models.Issue, issuesCount)
	for i := range issueslist {
		err := faker.FakeData(&issueslist[i])
		require.Nil(t, err)
	}

	t.Run("Create issue ok", func(t *testing.T) {
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(gomock.Any(), someUserOwner.Login, someRepo.Name).
				Return(true, nil).
				Times(1),
			repoIssues.EXPECT().
				CreateIssue(someIssue).
				Return(nil).
				Times(1),
		)

		err := useCase.CreateIssue(someIssue)

		require.NoError(t, err)
	})

	t.Run("Create issue error in repoID", func(t *testing.T) {
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, entityerrors.DoesNotExist()).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(0),
		)

		err := useCase.CreateIssue(someIssue)

		require.EqualError(t, err, entityerrors.DoesNotExist().Error())
	})

	t.Run("Create issue error in getLoginByID", func(t *testing.T) {
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, entityerrors.DoesNotExist()).
				Times(1),
		)

		err := useCase.CreateIssue(someIssue)

		require.Equal(t, err, entityerrors.DoesNotExist())
	})

	t.Run("Create issue error in checkReadAccess", func(t *testing.T) {
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(gomock.Any(), someUserOwner.Login, someRepo.Name).
				Return(true, errors.New("some error")).
				Times(1),
			repoIssues.EXPECT().
				CreateIssue(someIssue).
				Return(nil).
				Times(0),
		)

		err := useCase.CreateIssue(someIssue)

		require.Error(t, err)
	})

	t.Run("Create issue error in checkReadAccess", func(t *testing.T) {
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(gomock.Any(), someUserOwner.Login, someRepo.Name).
				Return(false, nil).
				Times(1),
		)

		err := useCase.CreateIssue(someIssue)

		require.Equal(t, err, entityerrors.AccessDenied())
	})

	t.Run("Update issues ok", func(t *testing.T) {
		gomock.InOrder(
			repoIssues.EXPECT().
				CheckEditAccessIssue(someIssue.AuthorID, someIssue.RepoID).
				Return(perm.WriteAccess(), nil).
				Times(1),
			repoIssues.EXPECT().
				UpdateIssue(someIssue).
				Return(nil).
				Times(1),
		)

		err := useCase.UpdateIssue(someIssue)

		require.NoError(t, err)
	})

	t.Run("Update issues access denied", func(t *testing.T) {
		gomock.InOrder(
			repoIssues.EXPECT().
				CheckEditAccessIssue(someIssue.AuthorID, someIssue.RepoID).
				Return(perm.NoAccess(), nil).
				Times(1),
			repoIssues.EXPECT().
				UpdateIssue(someIssue).
				Return(nil).
				Times(0),
		)

		err := useCase.UpdateIssue(someIssue)

		require.Equal(t, err, entityerrors.AccessDenied())
	})

	t.Run("GetIssue", func(t *testing.T) {
		gomock.InOrder(
			repoIssues.EXPECT().
				GetIssue(someIssue.ID).
				Return(someIssue, nil).
				Times(1),
		)

		issueFromDB, err := useCase.GetIssue(someIssue.ID, someUserOwner.ID)

		require.NoError(t, err)

		require.Equal(t, issueFromDB, someIssue)
	})

	t.Run("CloseIssue", func(t *testing.T) {
		gomock.InOrder(
			repoIssues.EXPECT().
				CheckEditAccessIssue(someIssue.AuthorID, someIssue.ID).
				Return(perm.AdminAccess(), nil).
				Times(1),
			repoIssues.EXPECT().
				CloseIssue(someIssue.ID).
				Return(nil).
				Times(1),
		)

		err := useCase.CloseIssue(someIssue.ID, someUserOwner.ID)

		require.NoError(t, err)
	})

	t.Run("CloseIssue no access", func(t *testing.T) {
		gomock.InOrder(
			repoIssues.EXPECT().
				CheckEditAccessIssue(someIssue.AuthorID, someIssue.ID).
				Return(perm.ReadAccess(), nil).
				Times(1),
			repoIssues.EXPECT().
				CloseIssue(someIssue.ID).
				Return(nil).
				Times(0),
		)

		err := useCase.CloseIssue(someIssue.ID, someUserOwner.ID)

		require.Error(t, err)
		require.Equal(t, err, entityerrors.AccessDenied())
	})

	t.Run("GetIssuesList ok", func(t *testing.T) {

		var limit int64 = 10
		var offset int64 = 2
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(gomock.Any(), someUserOwner.Login, someRepo.Name).
				Return(true, nil).
				Times(1),
			repoIssues.EXPECT().
				GetIssuesList(someRepo.ID, limit, offset).
				Return(issueslist, nil).
				Times(1),
		)

		issuesListFromDB, err := useCase.GetIssuesList(someRepo.ID, someUserOwner.ID, limit, offset)

		require.NoError(t, err)
		require.EqualValues(t, issueslist, issuesListFromDB)
	})
	t.Run("GetIssuesList err in get by id ", func(t *testing.T) {
		var limit int64 = 10
		var offset int64 = 2
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, errors.New("some error")).
				Times(1),
		)

		issuesListFromDB, err := useCase.GetIssuesList(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Error(t, err)
		require.EqualValues(t, models.IssuesSet{}, issuesListFromDB)
	})
	t.Run("GetIssuesList err in GetLoginByID ", func(t *testing.T) {
		var limit int64 = 10
		var offset int64 = 2
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, entityerrors.DoesNotExist()).
				Times(1),
		)

		issuesListFromDB, err := useCase.GetIssuesList(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Error(t, err)
		require.EqualValues(t, models.IssuesSet{}, issuesListFromDB)
	})

	t.Run("GetIssuesList err in checkReadAccess ", func(t *testing.T) {
		var limit int64 = 10
		var offset int64 = 2
		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(gomock.Any(), someUserOwner.Login, someRepo.Name).
				Return(false, nil).
				Times(1),
		)

		issuesListFromDB, err := useCase.GetIssuesList(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Error(t, err)
		require.EqualValues(t, models.IssuesSet{}, issuesListFromDB)

		gomock.InOrder(
			gitRepo.EXPECT().
				GetByID(someRepo.ID).
				Return(someRepo, nil).
				Times(1),
			userRepo.EXPECT().
				GetLoginByID(someUserOwner.ID).
				Return(someUserOwner.Login, nil).
				Times(1),
			gitRepo.EXPECT().
				CheckReadAccess(gomock.Any(), someUserOwner.Login, someRepo.Name).
				Return(false, errors.New("some error")).
				Times(1),
		)
		issuesListFromDB, err = useCase.GetIssuesList(someRepo.ID, someUserOwner.ID, limit, offset)

		require.Error(t, err)
		require.EqualValues(t, models.IssuesSet{}, issuesListFromDB)
	})

}
