package errs

import (
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
func (e *AppError) ToGRPC() error {
	switch e.Code {
	case http.StatusBadRequest:
		return status.Error(codes.InvalidArgument, e.Message)
	case http.StatusUnauthorized:
		return status.Error(codes.Unauthenticated, e.Message)
	case http.StatusNotFound:
		return status.Error(codes.NotFound, e.Message)
	case http.StatusConflict:
		return status.Error(codes.AlreadyExists, e.Message)
	case http.StatusUnprocessableEntity:
		return status.Error(codes.InvalidArgument, e.Message)
	default:
		return status.Error(codes.Internal, e.Message)
	}
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

func Unauthorized(err error) *AppError {
	return &AppError{
		Message: "unauthorized",
		Code:    http.StatusUnauthorized,
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
	case "categories_slug_key":
		msg = "slug already exists"
	case "products_slug_key":
		msg = "slug already exists"
	case "admin_users_email_key":
		msg = "email already exists"
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

func ForeignKeyViolation(err error, pgErr *pgconn.PgError) *AppError {
	var msg string
	switch pgErr.ConstraintName {
	case "fk_categories_parent_id":
		msg = "there is no category with such parent id"
	case "fk_products_brand_id":
		msg = "there is no brand with such id"
	case "fk_products_category_id":
		msg = "there is no category with such id"
	default:
		msg = "foreign key violation"
	}
	return &AppError{
		Message: msg,
		Code:    http.StatusUnprocessableEntity,
		Err:     err,
	}
}

func IsForeignKeyViolation(err error) (pgErr *pgconn.PgError, isUniqueViolation bool) {
	if err != nil && errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return pgErr, true
	}
	return nil, false
}
