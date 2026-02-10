package errs

import (
	"net/http"
)

type AppError struct {
	Message string
	Code    int
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func Internal(err error) *AppError {
	return &AppError{
		Message: "internal server error",
		Code:    http.StatusInternalServerError,
		Err:     err,
	}
}

func NotFound(err error) *AppError {
	return &AppError{
		Message: "not found",
		Code:    http.StatusNotFound,
		Err:     err,
	}
}

func BadRequest(err error) *AppError {
	return &AppError{
		Message: "bad request",
		Code:    http.StatusBadRequest,
		Err:     err,
	}
}
