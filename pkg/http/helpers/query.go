package helpers

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"strconv"
)

const (
	DefaultOffset int64 = 0
	DefaultLimit  int64 = 100
)

func ParseLimitAndOffset(values url.Values) (offset int64, limit int64, err error) {
	offset, err = strconv.ParseInt(values.Get("offset"), 10, 64)
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}
	if offset < 0 {
		return 0, 0, fmt.Errorf("offset parameter lower then zero")
	}

	limit, err = strconv.ParseInt(values.Get("limit"), 10, 64)
	if err != nil {
		return 0, 0, errors.WithStack(err)
	}
	if limit < 0 {
		return 0, 0, fmt.Errorf("limit parameter lower then zero")
	}

	return offset, limit, nil
}
