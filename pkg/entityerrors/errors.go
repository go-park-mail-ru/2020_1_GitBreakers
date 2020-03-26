package entityerrors

import "github.com/pkg/errors"

func DoesNotExist() error {
	return errors.New("entity does not exist")
}

func AlreadyExist() error {
	return errors.New("entity already not exist")
}

func Invalid() error {
	return errors.New("entity is invalid")
}
