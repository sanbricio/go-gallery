package imageRepository

import (
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
)

type ImageRepository interface {
	Find(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException)
	Insert(dto *imageDTO.ImageUploadRequestDTO) (*imageDTO.ImageDTO, *exception.ApiException)
	Update(dto *imageDTO.ImageUpdateRequestDTO) (*imageDTO.ImageUpdateResponseDTO, *exception.ApiException)
	Delete(dto *imageDTO.ImageDTO) (*imageDTO.ImageDTO, *exception.ApiException)
}
