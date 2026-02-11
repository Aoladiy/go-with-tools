package server

import (
	"go-with-tools/internal/errs"
	"strconv"

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

func getStringPathParam(c *gin.Context, param string) string {
	return c.Param(param)
}

func getIntPathParam(c *gin.Context, param string) (int, *errs.AppError) {
	pathParam, err := strconv.Atoi(c.Param(param))
	if err != nil {
		return 0, errs.BadRequest(err)
	}
	return pathParam, nil
}

func respondError(c *gin.Context, err *errs.AppError) {
	_ = c.Error(err)
	c.JSON(err.Code, gin.H{"error": err.Error()})
}
