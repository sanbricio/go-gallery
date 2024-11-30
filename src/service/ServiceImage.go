package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	repository "api-upload-photos/src/infrastructure/repository/image"
)

type ServiceImage struct {
	repository repository.ImageRepository
}

func NewServiceImage(repository repository.ImageRepository) *ServiceImage {
	return &ServiceImage{
		repository: repository,
	}
}

func (s *ServiceImage) Find(dtoFindImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Find(dtoFindImage)
}

func (s *ServiceImage) Insert(processedImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Insert(processedImage)
}

func (s *ServiceImage) Delete(id string) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Delete(id)
}
