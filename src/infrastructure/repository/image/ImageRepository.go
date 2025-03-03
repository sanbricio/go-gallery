package image_repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type ImageRepository interface {
	Insert(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Find(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Delete(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
}
