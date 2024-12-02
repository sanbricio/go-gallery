package image_repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type ImageRepository interface {
	Insert(dtoInsertImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Find(dtoFind *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Delete(id string) (*dto.DTOImage, *exception.ApiException)
}
