package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
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

func UniqueViolation(err error, pgErr *pgconn.PgError) *AppError {
	var msg string
	switch pgErr.ConstraintName {
	case "brands_name_key":
		msg = "name already exists"
	case "brands_slug_key":
		msg = "slug already exists"
	default:
		msg = "unique violation"
	}
	return &AppError{
		Message: msg,
		Code:    http.StatusConflict,
		Err:     err,
	}
}

func IsUniqueViolation(err error) (pgErr *pgconn.PgError, isUniqueViolation bool) {
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return pgErr, true
	}
	return nil, false
}
