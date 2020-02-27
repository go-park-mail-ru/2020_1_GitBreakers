package models

type statusMessageError struct {
	Code    int
	Message string
}

func (e *statusMessageError) GetCode() int {
	return e.Code
}

func (e *statusMessageError) Error() string {
	return e.Message
}

type CommonError interface {
	error
	GetCode() int
}

type ModelError struct {
	statusMessageError
}

func NewModelError(msg string, code int) *ModelError {
	err := new(ModelError)
	err.Message = msg
	err.Code = code
	return err
}