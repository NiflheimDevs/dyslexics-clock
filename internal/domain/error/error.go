package derror

import (
	"errors"
)

type ErrorType string

const (
	ErrTypeNotFound      ErrorType = "NOT_FOUND"
	ErrTypeUnprocessable ErrorType = "UNPROCESSABLE"
	ErrTypeUnauthorized  ErrorType = "UNAUTHORIZED"
	ErrTypeForbidden     ErrorType = "FORBIDDEN"
	ErrTypeBadRequest    ErrorType = "BAD_REQUEST"
	ErrTypeConflict      ErrorType = "CONFLICT"
	ErrTypeForeignKey    ErrorType = "FOREIGN_KEY"
	ErrTypeTimeout       ErrorType = "TIMEOUT"
	ErrTypeInternal      ErrorType = "INTERNAL"
	ErrTypeDB            ErrorType = "DB_ERROR"
)

type DomainError struct {
	Type ErrorType
	Err  error
	Msg  string
}

// func (e *DomainError) Error() string {
// 	if e.Err != nil {
// 		return fmt.Sprintf("%s: %s: %v", e.Type, e.Msg, e.Err)
// 	}
// 	return fmt.Sprintf("%s: %s", e.Type, e.Msg)
// }

func (e *DomainError) Error() string {
	return e.Msg
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

func New(kind ErrorType, msg string, err error) *DomainError {
	return &DomainError{
		Type: kind,
		Msg:  msg,
		Err:  err,
	}
}

func IsType(err error, kind ErrorType) bool {
	var e *DomainError
	if errors.As(err, &e) {
		return e.Type == kind
	}
	return false
}
