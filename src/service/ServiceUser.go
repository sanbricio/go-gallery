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

func (s *ServiceUser) Insert(dtoInsertUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Insert(dtoInsertUser)
}

func (s *ServiceUser) Find(dtoFindUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Find(dtoFindUser)
}

func (s *ServiceUser) FindAndCheckJWT(dtoFindUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.FindAndCheckJWT(dtoFindUser)
}

func (s *ServiceUser) Update(dtoUpdateUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Update(dtoUpdateUser)
}

func (s *ServiceUser) Delete(dtoDeleteUser *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return s.repository.Delete(dtoDeleteUser)
}
