package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Refresh receives a JWT and refreshes it
func Refresh(c *gin.Context) {
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

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		fmt.Println("Not enough time has elapsed for this token")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	// Create new token for user with a new Expiration Date
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{
		"name":    "token",
		"value":   tokenString,
		"expires": expirationTime,
	})
}
