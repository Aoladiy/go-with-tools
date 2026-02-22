package server

import (
	"context"
	"errors"
	"fmt"
	"go-with-tools/internal/auth"
	"go-with-tools/internal/errs"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
			respondError(c, errs.Unauthorized(errors.New(fmt.Sprintf("wrong token format (first word was not \"Bearer\". It was \"%s\")", bearerAndToken[0]))))
			c.Abort()
			return
		}
		token := bearerAndToken[1]
		parsedToken, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			if token.Method == jwt.SigningMethodHS256 {
				return []byte(jwtSecret), nil
			}
			return nil, fmt.Errorf("wrong signing method - %s", token.Method.Alg())
		})
		if err != nil {
			respondError(c, errs.Unauthorized(err))
			c.Abort()
			return
		}
		if !parsedToken.Valid {
			respondError(c, errs.Unauthorized(errors.New("token is invalid")))
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
		c.Set(auth.UserId, int64(userID))
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), auth.UserId, int64(userID)))
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
