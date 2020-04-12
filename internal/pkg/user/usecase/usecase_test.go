package usecase

import (
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

var someUser = models.User{
	ID:       12,
	Password: "sjfsfser242df",
	Name:     "Kekkers",
	Login:    "alahahbar",
	Image:    "/static/image/avatar/kek.png",
	Email:    "putin@kremlin.ru",
}

func TestUCUser_Create(t *testing.T) {

	t.Run("Create already exsist", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().IsExists(someUser).Return(true, nil)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Create(someUser)
		assert.Equal(t, err, entityerrors.AlreadyExist())
	})

}
func TestUCUser_Delete(t *testing.T) {

	t.Run("Delete ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().DeleteByID(someUser.ID).Return(nil)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Delete(someUser)
		assert.NoError(t, err)
	})
	t.Run("Delete not ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		someErr := errors.New("some error")

		m.EXPECT().DeleteByID(someUser.ID).Return(someErr)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Delete(someUser)
		assert.Equal(t, errors.Cause(err), someErr)
	})

}
func TestUCUser_CheckPass(t *testing.T) {

	t.Run("check pass false", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().CheckPass(someUser.Login, someUser.Password, ).Return(false, nil)

		useCase := UCUser{
			RepUser: m,
		}

		isCorrect, err := useCase.CheckPass(someUser.Login, someUser.Password)
		assert.Equal(t, false, isCorrect)
		assert.NoError(t, err)
	})
	t.Run("check pass true", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().CheckPass(someUser.Login, someUser.Password, ).
			Return(true, entityerrors.DoesNotExist())

		useCase := UCUser{
			RepUser: m,
		}

		isCorrect, err := useCase.CheckPass(someUser.Login, someUser.Password)
		assert.Equal(t, true, isCorrect)
		assert.Equal(t, err, entityerrors.DoesNotExist())
	})

}

func TestUCUser_GetByID(t *testing.T) {

	t.Run("get by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().GetUserByIDWithoutPass(someUser.ID).Return(someUser, nil)

		useCase := UCUser{
			RepUser: m,
		}

		userFromDB, err := useCase.GetByID(someUser.ID)
		assert.Equal(t, someUser, userFromDB)
		assert.NoError(t, err)
	})
}

func TestUCUser_GetByLogin(t *testing.T) {

	t.Run("get by login", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().GetByLoginWithoutPass(someUser.Login).Return(someUser, nil)

		useCase := UCUser{
			RepUser: m,
		}

		userFromDB, err := useCase.GetByLogin(someUser.Login)
		assert.Equal(t, someUser, userFromDB)
		assert.NoError(t, err)
	})
}

func TestUCUser_Update(t *testing.T) {

	t.Run("update with error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		m.EXPECT().GetUserByIDWithPass(someUser.ID).
			Return(someUser, entityerrors.DoesNotExist()).
			Times(1)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Update(someUser.ID, someUser)
		assert.Error(t, err)
	})
	t.Run("update ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		someUser.Password = "low"

		m.EXPECT().GetUserByIDWithPass(someUser.ID).
			Return(someUser, nil).
			Times(1)

		m.EXPECT().UserCanUpdate(someUser).
			Return(true, nil).
			Times(1)

		m.EXPECT().Update(someUser).
			Return(nil).
			Times(1)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Update(someUser.ID, someUser)
		assert.NoError(t, err)
	})
}

