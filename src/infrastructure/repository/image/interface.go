package imageRepository

import (
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
)

type ImageRepository interface {
	Insert(dto *imageDTO.ImageUploadRequestDTO, thumbnailImageId string) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException)
	Find(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException)
	Delete(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException)
}
