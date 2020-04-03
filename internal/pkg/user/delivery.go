package user

import (
	"io"
)

type UserDelivery interface {
	Update(writer io.Writer, reader *io.Reader)
	Login(writer io.Writer, reader *io.Reader)
	Create(writer io.Writer, reader *io.Reader)
	Logout(writer io.Writer, reader *io.Reader)
	Getinfo(writer io.Writer, reader *io.Reader)
	UploadAvatar(writer io.Writer, reader *io.Reader)
	GetInfoByLogin(writer io.Writer, reader *io.Reader)
}
