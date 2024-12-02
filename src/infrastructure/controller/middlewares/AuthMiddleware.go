package middlewares

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v4"
	jtoken "github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	secret  string
	handler fiber.Handler
}

func NewAuthMiddleware(secret string) *AuthMiddleware {
	return &AuthMiddleware{
		secret: secret,
		handler: jwtware.New(jwtware.Config{
			SigningKey: jwtware.SigningKey{Key: []byte(secret)},
		}),
	}
}

// Get values of JWT
func GetJWTClaims(token *jtoken.Token) (*dto.DTOUser, *exception.ApiException) {
	claims, ok := token.Claims.(jtoken.MapClaims)
	if !ok {
		return nil, exception.NewApiException(500, "invalid JWT claims type")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, exception.NewApiException(500, "username claim is missing or invalid")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, exception.NewApiException(500, "email claim is missing or invalid")
	}

	firstname, ok := claims["name"].(string)
	if !ok {
		return nil, exception.NewApiException(500, "name claim is missing or invalid")
	}

	return &dto.DTOUser{Username: username, Email: email, Firstname: firstname}, nil
}

func (auth *AuthMiddleware) GetSecret() string {
	return auth.secret
}

func (auth *AuthMiddleware) Handler() fiber.Handler {
	return auth.handler
}
