package serviceerror

import (
	"errors"
	"fmt"
	"net/http"
)

// Base error
var (
	ErrInvalidArgument  = fmt.Errorf("invalid argument")
	ErrNotFound         = fmt.Errorf("not found")
	ErrInternal         = fmt.Errorf("internal server error")
	ErrPermissionDenied = fmt.Errorf("permission denied")
	ErrUnauthenticated  = fmt.Errorf("unauthenticated")
	ErrConflict         = fmt.Errorf("conflict") // http 409 like
)

// Error is a custom error type, basically Base error but with custom info
type Error struct {
	StatusCode int
	Message    string
	Err        error
}

func NewError(wrappedErr error) *Error {
	return &Error{
		Err: wrappedErr,
	}
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) SetMessage(message string) *Error {
	e.Message = message

	return e
}

func (e *Error) SetStatusCode(status int) *Error {
	e.StatusCode = status

	return e
}

func NewInternal(err error) *Error {
	jErr := errors.Join(err, ErrInternal)

	return NewError(jErr).
		SetMessage("Internal server error.").
		SetStatusCode(http.StatusInternalServerError)
}

func NewNotFound(err error) *Error {
	jErr := errors.Join(err, ErrNotFound)

	return NewError(jErr).
		SetMessage("Not found.").
		SetStatusCode(http.StatusNotFound)
}

func NewInvalidArgument(err error) *Error {
	jErr := errors.Join(err, ErrInvalidArgument)

	return NewError(jErr).
		SetMessage("Invalid argument.").
		SetStatusCode(http.StatusBadRequest)
}

func NewUnauthenticated(err error) *Error {
	jErr := errors.Join(err, ErrUnauthenticated)

	return NewError(jErr).
		SetMessage("Unauthenticated.").
		SetStatusCode(http.StatusUnauthorized)
}

func NewPermissionDenied(err error) *Error {
	jErr := errors.Join(err, ErrPermissionDenied)

	return NewError(jErr).
		SetMessage("Permission denied.").
		SetStatusCode(http.StatusForbidden)
}
