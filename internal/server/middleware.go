package server

import "github.com/gin-gonic/gin"

func AuthByJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		//TODO actual functionality
		c.Next()
	}
}
