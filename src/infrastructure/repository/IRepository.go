package infrastructure

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
	handler "api-upload-photos/src/infrastructure"
)

type IRepository interface {
	Insert(fileInput *handler.ProcessedImage) (*dto.DTOImage, *exception.ApiException)
	Find(id string) (*dto.DTOImage, *exception.ApiException)
	Delete(id string) (*dto.DTOImage, *exception.ApiException)
}
