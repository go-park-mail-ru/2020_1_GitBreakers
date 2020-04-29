package usecase

import (
	mock_codehub "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/codehub/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	gitmodels "github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models/git"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
	m := mock_codehub.NewMockRepoStarI(ctrl)

	useCase := UCCodeHub{
		RepoStar: m,
	}

	t.Run("Modify star", func(t *testing.T) {

		m.EXPECT().AddStar(someUserOwner.ID, someRepo.ID).Return(nil)

		err := useCase.ModifyStar(models.Star{
			AuthorID: someUserOwner.ID,
			RepoID:   someRepo.ID,
			Vote:     true,
		})

		assert.NoError(t, err)
	})
}
