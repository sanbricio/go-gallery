package repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
)

type IRepositoryUser interface {
	Find(email, password string) (*entity.User, *exception.ApiException)
}
