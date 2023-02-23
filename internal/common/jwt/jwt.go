package jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Email    string `json:"email"`
	Id       string `json:"id"`
}

var (
	secret = []byte(os.Getenv("JWT_SECRET"))
)

// GenerateToken generate token and send it to user
func GenerateToken(claims Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(secret)
}
