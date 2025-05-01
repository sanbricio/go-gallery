package imageService

import (
	"go-gallery/src/commons/exception"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	imageRepository "go-gallery/src/infrastructure/repository/image"
	thumbnailImageRepository "go-gallery/src/infrastructure/repository/image/thumbnailImage"
)

type ImageService struct {
	imageRepository          imageRepository.ImageRepository
	thumbnailImageRepository thumbnailImageRepository.ThumbnailImageRepository
}

func NewImageService(imageRepository imageRepository.ImageRepository, thumbnailImageRepository thumbnailImageRepository.ThumbnailImageRepository) *ImageService {
	return &ImageService{
		imageRepository:          imageRepository,
		thumbnailImageRepository: thumbnailImageRepository,
	}
}

func (s *ImageService) Find(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	return s.imageRepository.Find(dto)
}

func (s *ImageService) Insert(dto *imageDTO.ImageUploadRequestDTO) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException) {
	thumbnailId, err := s.thumbnailImageRepository.Insert(dto)
	if err != nil {
		return nil, err
	}

	return s.imageRepository.Insert(dto, thumbnailId)
}

func (s *ImageService) Delete(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	return s.imageRepository.Delete(dto)
}
