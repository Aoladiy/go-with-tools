package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/Aoladiy/go-with-tools/internal/auth"
	"github.com/Aoladiy/go-with-tools/internal/config"
	"github.com/Aoladiy/go-with-tools/internal/errs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthByJWT(rdb *redis.Client, jwtSecret string) gin.HandlerFunc {
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
		isTokenSignedOut, err := auth.IsTokenSignedOut(c.Request.Context(), rdb, parsedToken.Raw)
		if err != nil {
			respondError(c, errs.Internal(err))
			c.Abort()
			return
		}
		if isTokenSignedOut {
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
		c.Next()
		if len(c.Errors) == 0 {
			return
		}

		log.Printf("\n[ERRORS %d] %s %s", len(c.Errors), c.Request.Method, c.FullPath())

		for i, e := range c.Errors {
			log.Printf("  Error %d:", i+1)
			logErrorRecursive(e.Err, "    ")
		}
		log.Printf("")
	}
}

func logErrorRecursive(err error, indent string) {
	if err == nil {
		return
	}

	log.Printf("%s -> %s", indent, err.Error())
	logErrorRecursive(errors.Unwrap(err), indent+"  ")
}
