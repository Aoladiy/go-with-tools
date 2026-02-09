package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomErr struct {
	Message string
	Code    int
}

func (e *CustomErr) Error() string {
	return e.Message
}

var (
	ErrNotFound = CustomErr{
		Message: "not found",
		Code:    http.StatusNotFound,
	}
	ErrBadRequest = CustomErr{
		Message: "bad request",
		Code:    http.StatusBadRequest,
	}
	ErrInternalServerError = CustomErr{
		Message: "Internal server error",
		Code:    http.StatusInternalServerError,
	}
)

func NonNil[T any](v []T) []T {
	if v == nil {
		return make([]T, 0)
	}

	return v
}

func ErrorResponse(c *gin.Context, err *CustomErr) {
	c.JSON(err.Code, gin.H{"error": err.Error()})
}
