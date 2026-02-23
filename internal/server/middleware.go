package server

import (
	"context"
	"errors"
	"go-with-tools/internal/auth"
	"go-with-tools/internal/config"
	"go-with-tools/internal/errs"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthByJWT(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		bearerAndToken := strings.Split(authorization, " ")
		if len(bearerAndToken) < 2 {
			respondError(c, errs.Unauthorized(errors.New("wrong token format")))
			c.Abort()
			return
		}
		if strings.ToLower(bearerAndToken[0]) != "bearer" {
			respondError(c, errs.Unauthorized(errors.New("invalid authorization header format")))
			c.Abort()
			return
		}
		token := bearerAndToken[1]
		parsedToken, appErr := auth.ParseToken(token, jwtSecret)
		if appErr != nil {
			respondError(c, appErr)
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
