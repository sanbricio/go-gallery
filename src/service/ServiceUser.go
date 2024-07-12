package service

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
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

func (s *ServiceUser) Find(username string, password string) (*entity.User, *exception.ApiException) {
	return s.repository.Find(username, password)
}
