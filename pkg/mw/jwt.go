package mw

import (
	"errors"
	"fmt"
	"grovo/internal/common/fiber/util"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JWTAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization", "")

	if !strings.Contains(authHeader, "Bearer") || authHeader == "" {
		return util.Send(c, 401, "invalid headers authorization")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", -1)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return c.Redirect("/login", 302)
		}
		return util.Send(c, 401, "invalid jsonwebtoken")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("token claims", err.Error())
		return util.Send(c, 401, "invalid jsonwebtoken")
	}

	c.Locals("jwt", claims)

	return c.Next()
}
