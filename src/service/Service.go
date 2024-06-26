package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
	handler "api-upload-photos/src/infrastructure"
	infrastructure "api-upload-photos/src/infrastructure/repository"
)

type Service struct {
	repository infrastructure.IRepository
}

func NewService(repository infrastructure.IRepository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Find(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Find(id)
}

func (s *Service) Insert(processedImage *handler.ProcessedImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Insert(processedImage)
}

func (s *Service) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Delete(id)
}
