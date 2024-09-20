package repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type IRepositoryImage interface {
	Insert(dtoInsertImage *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Find(dtoFind *dto.DTOImage) (*dto.DTOImage, *exception.ApiException)
	Delete(id string) (*dto.DTOImage, *exception.ApiException)
}
