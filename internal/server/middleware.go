package server

import (
	"errors"
	"go-with-tools/internal/errs"
	"log"

	"github.com/gin-gonic/gin"
)

func AuthByJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO actual functionality
		c.Next()
	}
}

func LogErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		for _, e := range c.Errors {
			var appErr *errs.AppError
			if errors.As(e.Err, &appErr) {
				if cause := appErr.Unwrap(); cause != nil {
					log.Printf("method \"%s\" | route \"%s\" | HTTP code \"%d\" | error \"%s\" | caused by \"%v\"", c.Request.Method, c.FullPath(), appErr.Code, appErr, cause)
				} else {
					log.Printf("method \"%s\" | route \"%s\" | HTTP code \"%d\" | error \"%s\"", c.Request.Method, c.FullPath(), appErr.Code, appErr)
				}
			} else {
				log.Printf("method \"%s\" | route \"%s\" | HTTP code \"%d\" | error \"%v\"", c.Request.Method, c.FullPath(), c.Writer.Status(), e.Error())
			}
		}
	}
}
