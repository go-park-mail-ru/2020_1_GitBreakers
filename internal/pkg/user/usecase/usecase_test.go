package usecase

import (
	"bytes"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/models"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/internal/pkg/user/mocks"
	"github.com/go-park-mail-ru/2020_1_GitBreakers/pkg/entityerrors"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"image"
	"image/color"
	"image/png"
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

func createFakeImage(width, height int) *image.RGBA {
	upLeft := image.Point{}
	lowRight := image.Point{X: width, Y: height}

	img := image.NewRGBA(image.Rectangle{Min: upLeft, Max: lowRight})

	// Colors are defined by Red, Green, Blue, Alpha uint8 values.
	cyan := color.RGBA{R: 100, G: 200, B: 200, A: 0xff}

	// Set color for each pixel.
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			switch {
			case x < width/2 && y < height/2: // upper left quadrant
				img.Set(x, y, cyan)
			case x >= width/2 && y >= height/2: // lower right quadrant
				img.Set(x, y, color.White)
			default:
				// Use zero value.
			}
		}
	}

	return img
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

	t.Run("Create some error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		someError := errors.New("some error")

		m.EXPECT().IsExists(someUser).Return(true, someError).Times(1)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Create(someUser)

		assert.Equal(t, someError, errors.Cause(err))
	})

	t.Run("Create good", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		gomock.InOrder(
			m.EXPECT().
				IsExists(someUser).
				Return(false, nil).
				Times(1),
			m.EXPECT().
				Create(gomock.Any()).
				Return(nil).
				Times(1),
		)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Create(someUser)

		require.Nil(t, err)
	})
	t.Run("Create err in creating", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		m := mocks.NewMockRepoUser(ctrl)

		someErr := errors.New("some error")

		gomock.InOrder(
			m.EXPECT().
				IsExists(someUser).
				Return(false, nil).
				Times(1),
			m.EXPECT().
				Create(gomock.Any()).
				Return(someErr).
				Times(1),
		)

		useCase := UCUser{
			RepUser: m,
		}

		err := useCase.Create(someUser)

		require.NotNil(t, err)
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

		m.EXPECT().CheckPass(someUser.Login, someUser.Password).Return(false, nil)

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

		m.EXPECT().CheckPass(someUser.Login, someUser.Password).
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepoUser(ctrl)
	useCase := UCUser{
		RepUser: m,
	}

	t.Run("update with error", func(t *testing.T) {
		m.EXPECT().GetUserByIDWithPass(someUser.ID).
			Return(someUser, entityerrors.DoesNotExist()).
			Times(1)

		err := useCase.Update(someUser.ID, someUser)
		assert.Error(t, err)
	})

	t.Run("update ok", func(t *testing.T) {

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

		err := useCase.Update(someUser.ID, someUser)
		assert.NoError(t, err)
	})

	t.Run("update full valid with conflict", func(t *testing.T) {

		m.EXPECT().GetUserByIDWithPass(someUser.ID).
			Return(someUser, nil).
			Times(1)

		m.EXPECT().UserCanUpdate(gomock.Any()).
			Return(false, nil).
			Times(1)

		m.EXPECT().Update(gomock.Any()).
			Return(nil).
			Times(0)

		err := useCase.Update(someUser.ID, someUser)

		require.Equal(t, entityerrors.AlreadyExist(), errors.Cause(err))
	})
	t.Run("update full ok", func(t *testing.T) {
		someUser.Password = "gooodpassword"

		m.EXPECT().GetUserByIDWithPass(someUser.ID).
			Return(someUser, nil).
			Times(1)

		m.EXPECT().UserCanUpdate(gomock.AssignableToTypeOf(someUser)).
			Return(true, nil).
			Times(1)

		m.EXPECT().Update(gomock.AssignableToTypeOf(someUser)).
			Return(nil).
			Times(1)

		err := useCase.Update(someUser.ID, someUser)
		assert.NoError(t, err)
	})
}
func TestUCUser_UploadAvatar(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockRepoUser(ctrl)

	someName := "some_name.png"

	binaryImage := new(bytes.Buffer)
	err := png.Encode(binaryImage, createFakeImage(10, 10))
	require.NoError(t, err)
	useCase := UCUser{m}

	t.Run("UploadAvatar ok", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				UploadAvatar(gomock.AssignableToTypeOf(someName), binaryImage.Bytes()).
				Return(nil).
				Times(1),
			m.EXPECT().
				GetUserByIDWithPass(someUser.ID).
				Return(someUser, nil).
				Times(1),
			m.EXPECT().
				UpdateAvatarPath(someUser, gomock.AssignableToTypeOf(someName)).
				Return(nil).
				Times(1),
		)

		err = useCase.UploadAvatar(someUser.ID, someName,
			binaryImage.Bytes())

		require.NoError(t, err)
	})

	t.Run("UploadAvatar incorrect content", func(t *testing.T) {
		someBadContent := []byte("fasfsffaf")
		gomock.InOrder(
			m.EXPECT().
				UploadAvatar(someName, binaryImage.Bytes()).
				Return(nil).
				Times(0),
		)

		err = useCase.UploadAvatar(someUser.ID, someName,
			someBadContent)

		require.Error(t, err)
	})

	t.Run("UploadAvatar err in uploadAvatar", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				UploadAvatar(gomock.AssignableToTypeOf(someName), binaryImage.Bytes()).
				Return(errors.New("some error")).
				Times(1),
		)

		err = useCase.UploadAvatar(someUser.ID, someName,
			binaryImage.Bytes())

		require.Error(t, err)
	})

	t.Run("UploadAvatar err in getByID", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				UploadAvatar(gomock.AssignableToTypeOf(someName), binaryImage.Bytes()).
				Return(nil).
				Times(1),
			m.EXPECT().
				GetUserByIDWithPass(someUser.ID).
				Return(someUser, errors.New("some error")).
				Times(1),
		)

		err = useCase.UploadAvatar(someUser.ID, someName,
			binaryImage.Bytes())

		require.Error(t, err)
	})
	t.Run("UploadAvatar err in updatePath", func(t *testing.T) {
		gomock.InOrder(
			m.EXPECT().
				UploadAvatar(gomock.AssignableToTypeOf(someName), binaryImage.Bytes()).
				Return(nil).
				Times(1),
			m.EXPECT().
				GetUserByIDWithPass(someUser.ID).
				Return(someUser, nil).
				Times(1),
			m.EXPECT().
				UpdateAvatarPath(someUser, gomock.AssignableToTypeOf(someName)).
				Return(errors.New("some error")).
				Times(1),
		)

		err = useCase.UploadAvatar(someUser.ID, someName,
			binaryImage.Bytes())

		require.Error(t, err)
	})

}
