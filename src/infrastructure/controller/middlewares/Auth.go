package middlewares

import (
	"api-upload-photos/src/infrastructure/dto"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v4"
	jtoken "github.com/golang-jwt/jwt/v5"
)

// Middleware JWT function
func NewAuthMiddleware(secret string) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(secret)},
	})
}

// Get values of JWT
func GetJWTClaims(token *jtoken.Token) *dto.DTOUser {
	claims := token.Claims.(jtoken.MapClaims)
	username := claims["username"].(string)
	email := claims["email"].(string)
	firstname := claims["name"].(string)

	return &dto.DTOUser{Username: username, Email: email, Firstname: firstname}
}
