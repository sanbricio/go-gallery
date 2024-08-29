package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	repository "api-upload-photos/src/infrastructure/repository/user"
)

type ServiceUser struct {
	repository repository.IRepositoryUser
}

func NewServiceUser(repository repository.IRepositoryUser) *ServiceUser {
	return &ServiceUser{
		repository: repository,
	}
}

func (s *ServiceUser) Find(dtoLoginRequest *dto.DTOLoginRequest) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Find(dtoLoginRequest)
}

func (s *ServiceUser) Insert(dtoRegisterRequest *dto.DTORegisterRequest) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Insert(dtoRegisterRequest)
}
