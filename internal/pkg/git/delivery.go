package git

import "io"

type GitDelivery interface {
	Create(writer io.Writer, reader *io.Reader)
}
