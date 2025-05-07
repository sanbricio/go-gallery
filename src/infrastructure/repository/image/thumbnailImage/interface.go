package thumbnailImageRepository

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/dto"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"
)

type ThumbnailImageRepository interface {
	Insert(dto *imageDTO.ImageDTO, rawContentFile []byte) (*imageDTO.ImageUploadResponseDTO, *exception.ApiException)
	Update(dto *imageDTO.ImageUpdateRequestDTO) (*imageDTO.ImageUpdateResponseDTO, *exception.ApiException)
	Delete(dto *imageDTO.ImageDeleteRequestDTO) (*dto.MessageResponseDTO, *exception.ApiException)
	FindAll(owner, lastIDHex string, pageSize int64) (*thumbnailImageDTO.ThumbnailImageCursorDTO, *exception.ApiException)
}
