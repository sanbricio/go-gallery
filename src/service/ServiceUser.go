package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	repository "api-upload-photos/src/infrastructure/repository/user"
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

func (s *UserService) Find(dtoFindUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Find(dtoFindUser)
}

func (s *UserService) FindAndCheckJWT(dtoFindUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.FindAndCheckJWT(dtoFindUser)
}

func (s *UserService) Update(dtoUpdateUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Update(dtoUpdateUser)
}

func (s *UserService) Delete(dtoDeleteUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Delete(dtoDeleteUser)
}
