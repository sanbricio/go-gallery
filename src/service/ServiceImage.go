package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	repository "api-upload-photos/src/infrastructure/repository/image"
)

type ServiceImage struct {
	repository repository.IRepositoryImage
}

func NewServiceImage(repository repository.IRepositoryImage) *ServiceImage {
	return &ServiceImage{
		repository: repository,
	}
}

func (s *ServiceImage) Find(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Find(id)
}

func (s *ServiceImage) Insert(processedImage *dto.DTOProcessedImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Insert(processedImage)
}

func (s *ServiceImage) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Delete(id)
}
