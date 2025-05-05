package userService

import (
	"go-gallery/src/commons/exception"

	userDTO "go-gallery/src/infrastructure/dto/user"
	repository "go-gallery/src/infrastructure/repository/user"
)

type UserService struct {
	repository repository.UserRepository
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Insert(userDTO *userDTO.UserDTO) (*userDTO.UserDTO, *exception.ApiException) {
	return s.repository.Insert(userDTO)
}

func (s *UserService) Find(loginRequestDTO *userDTO.LoginRequestDTO) (*userDTO.UserDTO, *exception.ApiException) {
	return s.repository.Find(loginRequestDTO)
}

func (s *UserService) FindAndCheckJWT(claimsDTO *userDTO.JwtClaimsDTO) (*userDTO.UserDTO, *exception.ApiException) {
	return s.repository.FindAndCheckJWT(claimsDTO)
}

func (s *UserService) Update(userDTO *userDTO.UserDTO) (int64, *exception.ApiException) {
	return s.repository.Update(userDTO)
}

func (s *UserService) Delete(userDTO *userDTO.UserDTO) (int64, *exception.ApiException) {
	return s.repository.Delete(userDTO)
}
