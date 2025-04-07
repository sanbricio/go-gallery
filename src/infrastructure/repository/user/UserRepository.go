package user_repository

import (
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/dto"
)

type UserRepository interface {
	Find(dtoLoginRequest *dto.DTOLoginRequest) (*dto.DTOUser, *exception.ApiException)
	FindAndCheckJWT(claims *dto.DTOClaimsJwt) (*dto.DTOUser, *exception.ApiException)
	Insert(dtoRegisterRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException)
	Update(dtoUserUpdate *dto.DTOUser) (int64, *exception.ApiException)
	Delete(dtoUserDelete *dto.DTOUser) (int64, *exception.ApiException)
}
