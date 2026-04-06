package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/Aoladiy/go-with-tools/gen"
	"github.com/Aoladiy/go-with-tools/internal/auth"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/errs"
	uuid2 "github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

func AuthByJWT(client gen.AuthMicroserviceClient, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, appErr := getJWTFromHeader(c)
		if appErr != nil {
			respondError(c, appErr)
			c.Abort()
			return
		}
		parsedToken, appErr := auth.ParseToken(token, jwtSecret)
		if appErr != nil {
			respondError(c, appErr)
			c.Abort()
			return
		}
		signedOutResponse, err := client.IsTokenSignedOut(c.Request.Context(), &gen.IsTokenSignedOutRequest{Token: parsedToken.Raw})
		if err != nil {
			respondError(c, errs.Internal(err))
			c.Abort()
			return
		}
		if signedOutResponse.IsTokenSignedOut {
			respondError(c, errs.Unauthorized(fmt.Errorf("token is in signed out tokens cache")))
			c.Abort()
			return
		}
		subject, err := parsedToken.Claims.GetSubject()
		if err != nil {
			respondError(c, errs.Unauthorized(errors.New("token isn't valid (no subject)")))
			c.Abort()
			return
		}
		userID, err := strconv.Atoi(subject)
		if err != nil {
			respondError(c, errs.Unauthorized(errors.New("token isn't valid (invalid user id)")))
			c.Abort()
			return
		}
		c.Set(config.UserIdKey, int64(userID))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), config.UserIdKey, int64(userID)))
		c.Next()
	}
}

func LogErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		uuid := uuid2.New().String()
		c.Set("request_id", uuid)

		start := time.Now()
		c.Next()
		finish := time.Since(start)
		if len(c.Errors) == 0 {
			return
		}
		for _, e := range c.Errors {
			var appErr *errs.AppError
			if !errors.As(e, &appErr) {
				logError(uuid, c.Request.Method, c.FullPath(), errs.Internal(fmt.Errorf("error is not of type *errs.AppError: %w", e.Err)))
				continue
			}
			if appErr.Code == errs.InternalErrCode {
				logError(uuid, c.Request.Method, c.FullPath(), appErr.Err)
				continue
			}
			logInfo(uuid, c.Request.Method, c.FullPath(), appErr.Err)
		}
		logDebug(uuid, c.Request.Method, c.FullPath(), finish)
	}
}

func logError(requestId, method, path string, err error) {
	slog.Error("failed request",
		"request_id", requestId,
		"method", method,
		"path", path,
		"err", err,
	)
}

func logInfo(requestId, method, path string, err error) {
	slog.Info("failed request",
		"request_id", requestId,
		"method", method,
		"path", path,
		"err", err,
	)
}

func logDebug(requestId, method, path string, duration time.Duration) {
	slog.Debug("request finished",
		"request_id", requestId,
		"method", method,
		"path", path,
		"duration", duration,
	)
}
