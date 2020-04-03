package entityerrors

import "github.com/pkg/errors"

type EntityError error

var (
	ErrDoesNotExist EntityError = errors.New("entity does not exist")
	ErrAlreadyExist EntityError = errors.New("entity already exist")
	ErrInvalid      EntityError = errors.New("entity is invalid")
	ErrAccessDenied EntityError = errors.New("access to entity denied")
)

func DoesNotExist() error {
	return ErrDoesNotExist
}

func AlreadyExist() error {
	return ErrAlreadyExist
}

func Invalid() error {
	return ErrInvalid
}

func AccessDenied() error {
	return ErrAccessDenied
}
