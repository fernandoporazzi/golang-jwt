package auth

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Signin handles new registrations
func Signin(c *gin.Context) {
	var creds Credentials

	// Get the JSON body and decode into creds
	err := c.Bind(&creds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	// Get the expected password from our in memory map
	expectedPwd, ok := users[creds.Username]
	if !ok || expectedPwd != creds.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// JWT expiration time is expressed in Unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
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
