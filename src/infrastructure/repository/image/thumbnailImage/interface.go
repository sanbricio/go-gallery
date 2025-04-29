package thumbnailImageRepository

import (
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
)

type ThumbnailImageRepository interface {
	Insert(dto *imageDTO.ImageUploadRequestDTO) (string, *exception.ApiException)
}
