package infrastructure

import (
	"api-upload-photos/src/domain"
	"mime/multipart"
)

type IRepository interface {
	Insert(fileInput *multipart.FileHeader) (*domain.Response, error)
	Find(id string) (*domain.DTOImage, error)
}
