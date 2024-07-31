package repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type IRepositoryUser interface {
	Find(email, password string) (*dto.DTOUser, *exception.ApiException)
	Insert() (*dto.DTOUser, *exception.ApiException)
	Update() (*dto.DTOUser, *exception.ApiException)
	Delete() (*dto.DTOUser, *exception.ApiException)
}
