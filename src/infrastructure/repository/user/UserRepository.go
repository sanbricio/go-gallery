package user_repository

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/dto"
)

type UserRepository interface {
	Find(dtoLoginRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	FindAndCheckJWT(dtoLoginRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Insert(dtoRegisterRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Update(dtoUserUpdate *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Delete(dtoUserDelete *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
}
