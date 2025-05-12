package userRepository

import (
	"go-gallery/src/commons/exception"
	userDTO "go-gallery/src/infrastructure/dto/user"
)

type UserRepository interface {
	Find(userLoginRequestDTO *userDTO.LoginRequestDTO) (*userDTO.UserDTO, *exception.ApiException)
	FindByEmail(email string) (*userDTO.UserDTO, *exception.ApiException)
	FindAndCheckJWT(claims *userDTO.JwtClaimsDTO) (*userDTO.UserDTO, *exception.ApiException)
	Insert(userDTO *userDTO.UserDTO) (*userDTO.UserDTO, *exception.ApiException)
	Update(userDTO *userDTO.UserDTO) (int64, *exception.ApiException)
	Delete(userDTO *userDTO.UserDTO) (int64, *exception.ApiException)
}
