package service

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
	repository "api-upload-photos/src/infrastructure/repository/image"
)

type ImageService struct {
	repository repository.ImageRepository
}

func NewImageService(repository repository.ImageRepository) *ImageService {
	return &ImageService{
		repository: repository,
	}
}

func (s *ImageService) Find(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Find(dto)
}

func (s *ImageService) Insert(processedImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Insert(processedImage)
}

func (s *ImageService) Delete(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException) {
	return s.repository.Delete(dto)
}
