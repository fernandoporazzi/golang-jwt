package main

import (
	"github.com/fernandoporazzi/golang-jwt/auth"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	r.POST("/signin", auth.Signin)
	r.Run()
}
