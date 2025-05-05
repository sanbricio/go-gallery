package imageService

import (
	"go-gallery/src/commons/exception"

	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"
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
	imageDTO, err := s.imageRepository.Insert(dto)
	if err != nil {
		return nil, err
	}

	return s.thumbnailImageRepository.Insert(imageDTO, dto.RawContentFile)
}

func (s *ImageService) Update(dto *imageDTO.ImageUpdateRequestDTO) (*imageDTO.ImageUpdateResponseDTO, *exception.ApiException) {
	//TODO Añadir update del thumbnail
	return s.imageRepository.Update(dto)
}

func (s *ImageService) Delete(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException) {
	//TODO Añadir delete del thumbnail
	return s.imageRepository.Delete(dto)
}

func (s *ImageService) FindAllThumbnails(owner, lastID string, pageSize int64) (*thumbnailImageDTO.ThumbnailImageCursorDTO, *exception.ApiException) {
	return s.thumbnailImageRepository.FindAll(owner, lastID, pageSize)
}
