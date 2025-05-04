package thumbnailImageRepository

import (
	"go-gallery/src/commons/exception"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	thumbnailImageDTO "go-gallery/src/infrastructure/dto/image/thumbnailImage"
)

type ThumbnailImageRepository interface {
	Insert(dto *imageDTO.ImageUploadRequestDTO) (string, *exception.ApiException)
	FindAll(owner, lastIDHex string, pageSize int64) (*thumbnailImageDTO.ThumbnailImageCursorDTO, *exception.ApiException)
}
