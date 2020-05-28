package entityerrors

import "github.com/pkg/errors"

type EntityError error

var (
	ErrDoesNotExist EntityError = errors.New("entity does not exist")
	ErrAlreadyExist EntityError = errors.New("entity already exist")
	ErrInvalid      EntityError = errors.New("entity is invalid")
	ErrAccessDenied EntityError = errors.New("access to entity denied")
	ErrContentEmpty EntityError = errors.New("entity content empty")
	ErrConflict     EntityError = errors.New("entity conflicts with other entities")
	ErrTooLarge     EntityError = errors.New("entity too large")
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

func ContentEmpty() error {
	return ErrContentEmpty
}

func Conflict() error {
	return ErrConflict
}

func TooLarge() error {
	return ErrTooLarge
}
