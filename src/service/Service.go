package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
	entity "api-upload-photos/src/domain/entities"
	infrastructure "api-upload-photos/src/infrastructure/repository"
	"mime/multipart"
)

type Service struct {
	repository infrastructure.IRepository
}

func NewService(repository infrastructure.IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Insert(fileInput *multipart.FileHeader) (*entity.Response, *exception.ApiException) {
	return s.repository.Insert(fileInput)
}

func (s *Service) Find(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Find(id)
}
