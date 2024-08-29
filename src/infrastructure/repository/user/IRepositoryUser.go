package repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type IRepositoryUser interface {
	Find(dtoLoginRequest *dto.DTOLoginRequest) (*dto.DTOUser, *exception.ApiException)
	Insert(dtoRegister *dto.DTORegisterRequest) (*dto.DTOUser, *exception.ApiException)
	Update() (*dto.DTOUser, *exception.ApiException)
	Delete() (*dto.DTOUser, *exception.ApiException)
}
