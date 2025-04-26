package imageService

import (
	"go-gallery/src/commons/exception"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	repository "go-gallery/src/infrastructure/repository/image"
)

type ImageService struct {
	repository repository.ImageRepository
}

func NewImageService(repository repository.ImageRepository) *ImageService {
	return &ImageService{
		repository: repository,
	}
}

func (s *ImageService) Find(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	return s.repository.Find(dto)
}

func (s *ImageService) Insert(processedImage *imageDTO.ImageUploadRequestDTO) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException) {
	return s.repository.Insert(processedImage)
}

func (s *ImageService) Delete(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	return s.repository.Delete(dto)
}
