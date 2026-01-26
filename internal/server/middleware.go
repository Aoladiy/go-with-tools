package server

import "github.com/gin-gonic/gin"

func authByJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO actual functionality
		c.Next()
	}
}
