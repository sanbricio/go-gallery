package auth

import (
	"go-gallery/src/commons/exception"
	userDTO "go-gallery/src/infrastructure/dto/user"
)

type TokenManager interface {
	CreateToken(username, email string) (string, *exception.ApiException)
	ValidateToken(tokenString string) (*userDTO.JwtClaimsDTO, *exception.ApiException)
}
