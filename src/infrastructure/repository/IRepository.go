package infrastructure

import (
	"api-upload-photos/src/domain/dto"
	entity "api-upload-photos/src/domain/entities"
	"mime/multipart"
)

type IRepository interface {
	Insert(fileInput *multipart.FileHeader) (*entity.Response, error)
	Find(id string) (*dto.DTOImage, error)
}
