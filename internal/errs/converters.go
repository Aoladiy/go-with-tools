package errs

import (
	"net/http"
)

func (e *AppError) HttpCode() int {
	switch e.Code {
	case UnauthorizedErrCode:
		return http.StatusUnauthorized
	case BadRequestErrCode:
		return http.StatusBadRequest
	case ConflictErrCode:
		return http.StatusConflict
	case NotFoundErrCode:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
