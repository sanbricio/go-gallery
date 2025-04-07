package service

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/dto"
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

func (s *UserService) Insert(dtoInsertUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Insert(dtoInsertUser)
}

func (s *UserService) Find(dtoFindUser *dto.DTOLoginRequest) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Find(dtoFindUser)
}

func (s *UserService) FindAndCheckJWT(claims *dto.DTOClaimsJwt) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.FindAndCheckJWT(claims)
}

func (s *UserService) Update(dtoUpdateUser *dto.DTOUser) (int64, *exception.ApiException) {
	return s.repository.Update(dtoUpdateUser)
}

func (s *UserService) Delete(dtoDeleteUser *dto.DTOUser) (int64, *exception.ApiException) {
	return s.repository.Delete(dtoDeleteUser)
}
