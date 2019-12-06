package auth

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Me hhhahah
func Me(c *gin.Context) {
	authorizationHeader := c.Request.Header.Get("Authorization")

	if !strings.HasPrefix(strings.ToLower(authorizationHeader), "bearer") {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	tokenString := authorizationHeader[7:]
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			c.JSON(http.StatusUnauthorized, gin.H{})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	if !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}

	c.JSON(200, gin.H{
		"message": "hi " + claims.Username,
	})
}
