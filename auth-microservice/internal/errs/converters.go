package errs

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func (e *AppError) HttpCode() int {
	switch e.Code {
	case UnauthorizedErrCode:
		return http.StatusUnauthorized
	case BadRequestErrCode:
		return http.StatusBadRequest
	case UnprocessableEntityErrCode:
		return http.StatusUnprocessableEntity
	case ConflictErrCode:
		return http.StatusConflict
	case NotFoundErrCode:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func (e *AppError) GrpcCode() codes.Code {
	switch e.Code {
	case UnauthorizedErrCode:
		return codes.Unauthenticated
	case BadRequestErrCode:
		return codes.InvalidArgument
	case UnprocessableEntityErrCode:
		return codes.InvalidArgument
	case ConflictErrCode:
		return codes.AlreadyExists
	case NotFoundErrCode:
		return codes.NotFound
	default:
		return codes.Internal
	}
}