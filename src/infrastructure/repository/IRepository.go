package infrastructure

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
	entity "api-upload-photos/src/domain/entities"
	"mime/multipart"
)

type IRepository interface {
	Insert(fileInput *multipart.FileHeader) (*entity.Response, *exception.ApiException)
	Find(id string) (*dto.DTOImage, *exception.ApiException)
	Delete(id string) (*dto.DTOImage, *exception.ApiException)
}
