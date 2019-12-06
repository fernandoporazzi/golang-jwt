package auth

import "github.com/dgrijalva/jwt-go"

// Credentials handles user inputs
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Claims is the payload
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
