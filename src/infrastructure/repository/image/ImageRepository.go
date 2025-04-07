package image_repository

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/dto"
)

type ImageRepository interface {
	Insert(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Find(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Delete(dto *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
}
