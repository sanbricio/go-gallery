package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	repository "api-upload-photos/src/infrastructure/repository/image"
)

type Service struct {
	repository repository.IRepositoryImage
}

func NewService(repository repository.IRepositoryImage) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Find(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Find(id)
}

func (s *Service) Insert(processedImage *dto.ProcessedImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Insert(processedImage)
}

func (s *Service) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Delete(id)
}
