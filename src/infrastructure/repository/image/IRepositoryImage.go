package repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type IRepositoryImage interface {
	Insert(fileInput *dto.ProcessedImage) (*dto.DTOImage, *exception.ApiException)
	Find(id string) (*dto.DTOImage, *exception.ApiException)
	Delete(id string) (*dto.DTOImage, *exception.ApiException)
}
