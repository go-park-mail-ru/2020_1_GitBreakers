package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/git/mocks"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
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

func TestGitUseCase_GetRepo(t *testing.T) {

	t.Run("Get repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		username := "keker"
		repoName := "mdasher"
		userid := 5

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
}

func TestGitUseCase_Create(t *testing.T) {

	t.Run("Get repo", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepository(ctrl)
		//username := "keker"
		//repoName := "mdasher"
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
}
