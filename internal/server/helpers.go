package server

import (
	"errors"
	"strconv"
	"strings"

	"github.com/Aoladiy/go-with-tools/internal/DTO"
	"github.com/Aoladiy/go-with-tools/internal/errs"

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

func getInt64PathParam(c *gin.Context, param string) (int64, *errs.AppError) {
	pathParam, err := strconv.ParseInt(c.Param(param), 10, 64)
	if err != nil {
		return 0, errs.BadRequest(err)
	}
	return pathParam, nil
}

func respondError(c *gin.Context, err *errs.AppError) {
	_ = c.Error(err)
	c.JSON(
		err.Code,
		DTO.ErrorResponse{Error: err.Error()},
	)
}

func getJWTFromHeader(c *gin.Context) (string, *errs.AppError) {
	authorization := c.GetHeader("Authorization")
	bearerAndToken := strings.Split(authorization, " ")
	if len(bearerAndToken) < 2 || strings.TrimSpace(bearerAndToken[1]) == "" {
		return "", errs.Unauthorized(errors.New("wrong token format"))
	}
	if strings.ToLower(bearerAndToken[0]) != "bearer" {
		return "", errs.Unauthorized(errors.New("invalid authorization header format"))
	}
	return bearerAndToken[1], nil
}
