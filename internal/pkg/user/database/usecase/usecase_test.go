package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var someUser = models.User{
	Id:       12,
	Password: "sjfsfser242df",
	Name:     "Kekkers",
	Login:    "alahahbar",
	Image:    "/static/image/avatar/kek.png",
	Email:    "putin@kremlin.ru",
}

func TestCreate(t *testing.T) {

	t.Run("Create already exsist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepoUser(ctrl)

		m.EXPECT().IsExists(someUser).Return(true, nil)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Create(someUser)
		assert.Equal(t, err, entityerrors.AlreadyExist())
	})

}
func TestDelete(t *testing.T) {

	t.Run("Delete ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepoUser(ctrl)

		m.EXPECT().DeleteById(someUser.Id).Return(nil)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Delete(someUser)
		assert.NoError(t, err)
	})
	t.Run("Delete not ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := user.NewMockRepoUser(ctrl)

		someErr := errors.New("some error")

		m.EXPECT().DeleteById(someUser.Id).Return(someErr)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Delete(someUser)
		assert.Equal(t, errors.Cause(err), someErr)
	})

}
