package server

import (
	"errors"
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
