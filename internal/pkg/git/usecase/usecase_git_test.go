package usecase

import (
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

var someRepo = gitmodels.Repository{
	ID:          50,
	OwnerID:     12,
	Name:        "PythonProject",
	Description: "repo for work",
	IsFork:      false,
	CreatedAt:   time.Now(),
	IsPublic:    true,
}
var someUser = models.User{
	ID:       12,
	Password: "2392g3r39rri3ty33FSG3",
	Name:     "ssgsgsgsdg",
	Login:    "tiktak",
	Image:    "standart.png",
	Email:    "keksik@yandex.mda",
}

func TestGitUseCase_GetRepo(t *testing.T) {
	username := "keker"
	repoName := "mdasher"
	userid := 5

	t.Run("Get repo ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		m.EXPECT().
			CheckReadAccess(&userid, username, repoName).
			Return(true, nil)

		m.EXPECT().
			GetByName(username, repoName).
			Return(someRepo, nil)

		useCase := GitUseCase{
			Repo: m,
		}

		repoFromDb, err := useCase.GetRepo(username, repoName, &userid)
		require.Nil(t, err)
		require.Equal(t, repoFromDb, someRepo)
	})

	t.Run("Get repo access denied", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		m.EXPECT().
			CheckReadAccess(&userid, username, repoName).
			Return(false, nil)

		useCase := GitUseCase{
			Repo: m,
		}

		repoFromDb, err := useCase.GetRepo(username, repoName, &userid)

		require.Equal(t, repoFromDb, gitmodels.Repository{})
		require.Equal(t, err, entityerrors.AccessDenied())
	})

	t.Run("Get repo some error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		m.EXPECT().
			CheckReadAccess(&userid, username, repoName).
			Return(false, errors.New("some error"))

		useCase := GitUseCase{
			Repo: m,
		}

		repoFromDb, err := useCase.GetRepo(username, repoName, &userid)

		require.Equal(t, repoFromDb, gitmodels.Repository{})
		require.Error(t, err)
	})
}

func TestGitUseCase_Create(t *testing.T) {

	t.Run("Create repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		userid := 12

		m.EXPECT().
			Create(someRepo).
			Return(int64(45242), nil)

		useCase := GitUseCase{
			Repo: m,
		}

		err := useCase.Create(userid, &someRepo)
		require.Nil(t, err)

	})

	t.Run("Create repo already exsist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		userid := 12

		m.EXPECT().
			Create(someRepo).
			Return(int64(45242), entityerrors.AlreadyExist())

		useCase := GitUseCase{
			Repo: m,
		}

		err := useCase.Create(userid, &someRepo)
		require.Equal(t, err, entityerrors.AlreadyExist())

	})

	t.Run("Create repo some error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		userid := 12

		m.EXPECT().
			Create(someRepo).
			Return(int64(45242), errors.New("some error"))

		useCase := GitUseCase{
			Repo: m,
		}

		err := useCase.Create(userid, &someRepo)
		require.Error(t, err)

	})
}
func TestGitUseCase_GetRepoList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepository(ctrl)

	useCase := GitUseCase{
		Repo: m,
	}

	const repoCount int = 7
	repolist := make([]gitmodels.Repository, repoCount)
	for i := range repolist {
		err := faker.FakeData(&repolist[i])
		require.Nil(t, err)
	}

	t.Run("Get repo list ok", func(t *testing.T) {
		gomock.InOrder(m.EXPECT().
			GetAnyReposByUserLogin(someUser.Login, 0, 100).
			Return(repolist, nil).Times(1),
			m.EXPECT().
				CheckReadAccess(&someUser.ID, gomock.Any(), gomock.Any()).
				Return(true, nil).Times(repoCount))

		repolistFromDB, err := useCase.GetRepoList(someUser.Login, &someUser.ID)

		require.Nil(t, err)
		require.Equal(t, repolist, repolistFromDB)
	})

	t.Run("Get repo list some err", func(t *testing.T) {
		gomock.InOrder(m.EXPECT().
			GetAnyReposByUserLogin(someUser.Login, 0, 100).
			Return(repolist, errors.New("some error")).Times(1),
			m.EXPECT().
				CheckReadAccess(&someUser.ID, gomock.Any(), gomock.Any()).
				Return(true, nil).Times(0))

		repolistFromDB, err := useCase.GetRepoList(someUser.Login, &someUser.ID)

		require.Error(t, err)
		require.Nil(t, repolistFromDB)
	})
}

func TestGitUseCase_GetBranchList(t *testing.T) {
	const branchCount int = 5

	branchlist := make([]gitmodels.Branch, branchCount)

	for i := range branchlist {
		err := faker.FakeData(&branchlist[i])
		require.Nil(t, err)
	}

	t.Run("Get branch list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, someUser.Login, someRepo.Name).
				Return(true, nil).Times(1),
			m.EXPECT().
				GetBranchesByName(someUser.Login, someRepo.Name).
				Return(branchlist, nil).Times(1))

		useCase := GitUseCase{
			Repo: m,
		}

		repolistFromDB, err := useCase.GetBranchList(&someUser.ID, someUser.Login, someRepo.Name)
		require.Nil(t, err)
		require.Equal(t, branchlist, repolistFromDB)

	})

	t.Run("Not get branch list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, someUser.Login, someRepo.Name).
				Return(false, nil).Times(1),
			m.EXPECT().
				GetBranchesByName(someUser.Login, someRepo.Name).
				Return(branchlist, nil).Times(0))

		useCase := GitUseCase{
			Repo: m,
		}

		repolistFromDB, err := useCase.GetBranchList(&someUser.ID, someUser.Login, someRepo.Name)

		require.Equal(t, err, entityerrors.AccessDenied())

		require.Nil(t, repolistFromDB)

	})
}

func TestGitUseCase_FilesInCommitByPath(t *testing.T) {
	const filesCount int = 5

	fileslist := make([]gitmodels.FileInCommit, filesCount)

	for i := range fileslist {
		err := faker.FakeData(&fileslist[i])
		require.Nil(t, err)
	}
	commitRequest := gitmodels.FilesCommitRequest{
		UserName:    "keksik",
		Reponame:    "tortik",
		HashCommits: "gerkti3592go2290244g353",
		Path:        "/",
	}

	t.Run("Get files list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, commitRequest.UserName, commitRequest.Reponame).
				Return(true, nil).Times(1),
			m.EXPECT().
				FilesInCommitByPath(commitRequest.UserName, commitRequest.Reponame, commitRequest.HashCommits, commitRequest.Path).
				Return(fileslist, nil).Times(1))

		useCase := GitUseCase{
			Repo: m,
		}

		fileslistFromDB, err := useCase.FilesInCommitByPath(commitRequest, &someUser.ID)

		require.Nil(t, err)

		require.Equal(t, fileslist, fileslistFromDB)

	})
	t.Run("Get null files list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, commitRequest.UserName, commitRequest.Reponame).
				Return(false, nil).Times(1),
			m.EXPECT().
				FilesInCommitByPath(commitRequest.UserName, commitRequest.Reponame, commitRequest.HashCommits, commitRequest.Path).
				Return(fileslist, nil).Times(0))

		useCase := GitUseCase{
			Repo: m,
		}

		fileslistFromDB, err := useCase.FilesInCommitByPath(commitRequest, &someUser.ID)

		require.Nil(t, fileslistFromDB)

		require.Equal(t, err, entityerrors.AccessDenied())

	})

}
func TestGitUseCase_GetCommitsByCommitHash(t *testing.T) {
	const commitsCount int = 5

	commitslist := make([]gitmodels.Commit, commitsCount)

	for i := range commitslist {
		err := faker.FakeData(&commitslist[i])
		require.Nil(t, err)
	}
	commitRequest := gitmodels.CommitRequest{
		UserLogin:  "keksik",
		RepoName:   "batya",
		CommitHash: "fwfw5290rf3024",
		Offset:     0,
		Limit:      100,
	}

	t.Run("Get commits by commit hash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, commitRequest.UserLogin, commitRequest.RepoName).
				Return(true, nil).Times(1),
			m.EXPECT().
				GetCommitsByCommitHash(commitRequest.UserLogin,
					commitRequest.RepoName, commitRequest.CommitHash,
					commitRequest.Offset, commitRequest.Limit).
				Return(commitslist, nil).Times(1))

		useCase := GitUseCase{
			Repo: m,
		}

		commitslistFromDB, err := useCase.GetCommitsByCommitHash(commitRequest, &someUser.ID)

		require.Nil(t, err)

		require.Equal(t, commitslist, commitslistFromDB)

	})

	t.Run("Get commits by commit hash", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, commitRequest.UserLogin, commitRequest.RepoName).
				Return(false, nil).Times(1),
			m.EXPECT().
				GetCommitsByCommitHash(commitRequest.UserLogin,
					commitRequest.RepoName, commitRequest.CommitHash,
					commitRequest.Offset, commitRequest.Limit).
				Return(commitslist, nil).Times(0))

		useCase := GitUseCase{
			Repo: m,
		}

		commitslistFromDB, err := useCase.GetCommitsByCommitHash(commitRequest, &someUser.ID)

		require.Nil(t, commitslistFromDB)

		require.Equal(t, err, entityerrors.AccessDenied())

	})
}
func TestGitUseCase_GetCommitsByBranchName(t *testing.T) {
	const commitsCount int = 5

	commitslist := make([]gitmodels.Commit, commitsCount)

	for i := range commitslist {
		err := faker.FakeData(&commitslist[i])
		require.Nil(t, err)
	}

	t.Run("Get commits by branchname ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		branchName := "deploy"
		offset := 5
		limit := 40
		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, someUser.Login, someRepo.Name).
				Return(true, nil).Times(1),
			m.EXPECT().
				GetCommitsByBranchName(someUser.Login, someRepo.Name, branchName, offset, limit).
				Return(commitslist, nil).Times(1))

		useCase := GitUseCase{
			Repo: m,
		}

		commitslistFromDB, err := useCase.
			GetCommitsByBranchName(someUser.Login, someRepo.Name, branchName, offset, limit, &someUser.ID)

		require.Nil(t, err)

		require.Equal(t, commitslist, commitslistFromDB)

	})
	t.Run("Get commits by branchname access denied", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		branchName := "deploy"
		offset := 5
		limit := 40
		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, someUser.Login, someRepo.Name).
				Return(false, nil).Times(1),
			m.EXPECT().
				GetCommitsByBranchName(someUser.Login, someRepo.Name, branchName, offset, limit).
				Return(commitslist, nil).Times(0))

		useCase := GitUseCase{
			Repo: m,
		}

		commitslistFromDB, err := useCase.
			GetCommitsByBranchName(someUser.Login, someRepo.Name, branchName, offset, limit, &someUser.ID)

		require.Nil(t, commitslistFromDB)

		require.Equal(t, err, entityerrors.AccessDenied())
	})

}
func TestGitUseCase_GetFileByPath(t *testing.T) {
	request := gitmodels.FilesCommitRequest{
		UserName:    "kekers",
		Reponame:    "mdasher",
		HashCommits: "fsf3495gk394jg335",
		Path:        "readme.md",
	}
	retFile := gitmodels.FileCommitted{
		FileInfo: gitmodels.FileInCommit{
			Name:        "tikakfaf",
			FileType:    "blob",
			FileMode:    "blob",
			FileSize:    50,
			IsBinary:    false,
			ContentType: "application/json",
			EntryHash:   "vefk30242kf2942",
		},
		Content: `{"hello guys"}`,
	}

	t.Run("Get file access denied", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, request.UserName, request.Reponame).
				Return(false, nil).Times(1),
			m.EXPECT().
				GetFileByPath(request.UserName, request.Reponame, request.HashCommits, request.Path).
				Return(retFile, nil).Times(0))

		useCase := GitUseCase{
			Repo: m,
		}

		fileFromDB, err := useCase.GetFileByPath(request, &someUser.ID)

		require.Equal(t, err, entityerrors.AccessDenied())

		require.Equal(t, fileFromDB, gitmodels.FileCommitted{})
	})
	t.Run("Get file ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)

		gomock.InOrder(
			m.EXPECT().
				CheckReadAccess(&someUser.ID, request.UserName, request.Reponame).
				Return(true, nil).Times(1),
			m.EXPECT().
				GetFileByPath(request.UserName, request.Reponame, request.HashCommits, request.Path).
				Return(retFile, nil).Times(1))

		useCase := GitUseCase{
			Repo: m,
		}

		fileFromDB, err := useCase.GetFileByPath(request, &someUser.ID)

		require.Equal(t, err, nil)

		require.Equal(t, fileFromDB, retFile)
	})
}
