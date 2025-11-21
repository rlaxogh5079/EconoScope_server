package errs

import (
	"errors"
	"net/http"
)

type AppError struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code string, message string, status int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
	}
}

func Wrap(code string, message string, status int, err error) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: status,
		Err:        err,
	}
}

func FromError(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	return Wrap("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError, err)
}

var (
	ErrBadRequest = New("BAD_REQUEST", "bad request", http.StatusBadRequest)

	ErrUnauthorized = New("UNAUTHORIZED", "unauthorized", http.StatusUnauthorized)

	ErrForbidden = New("FORBIDDEN", "forbidden", http.StatusForbidden)

	ErrNotFound = New("NOT_FOUND", "resource not found", http.StatusNotFound)

	ErrInternal = New("INTERNAL_ERROR", "internal server error", http.StatusInternalServerError)
)
