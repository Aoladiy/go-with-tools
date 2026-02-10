package server

import (
	"go-with-tools/internal/errs"
	"log"

	"github.com/gin-gonic/gin"
)

func nonNilSlice[T any](v []T) []T {
	if v == nil {
		return make([]T, 0)
	}

	return v
}

func bindJson[T any](c *gin.Context) (T, *errs.AppError) {
	var request T
	if err := c.ShouldBindJSON(&request); err != nil {
		return request, errs.BadRequest(err)
	}
	return request, nil
}

func errorResponse(c *gin.Context, err *errs.AppError) {
	c.JSON(err.Code, gin.H{"error": err.Error()})
}

func fail(c *gin.Context, message string, appError *errs.AppError) {
	if cause := appError.Unwrap(); cause != nil {
		log.Printf("%s %s %s (is caused by %v)", c.Request.Method, c.FullPath(), message, cause)
	} else {
		log.Printf("%s %s %s (error %v)", c.Request.Method, c.FullPath(), message, appError)
	}
	errorResponse(c, appError)
}
