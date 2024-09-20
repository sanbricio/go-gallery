package repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type IRepositoryUser interface {
	Find(dtoLoginRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	FindJWT(dtoLoginRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Insert(dtoRegisterRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Update(dtoUserUpdate *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Delete(dtoUserDelete *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
}
