package user_repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"api-upload-photos/src/infrastructure/dto"
	"os"
)

const UserMemoryRepositoryKey = "UserMemoryRepository"

type UserMemoryRepository struct {
}

func NewUserMemoryRepository() UserRepository {
	return new(UserMemoryRepository)
}

func (r *UserMemoryRepository) Find(dtoLoginRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {

	files, err := os.ReadDir("data")
	if err != nil {
		return nil, exception.NewApiException(500, "Error al leer el directorio")
	}

	for _, file := range files {
		if file.Name() == "prueba.json" {
			if dtoLoginRequest.Username == "test" && dtoLoginRequest.Password == "test12345" {
				entity := entity.NewUser("test", "test12345", "test@mail.com", "Prueba", "Prueba2")
				dto := dto.FromUser(entity)
				return dto, nil
			}
		}
	}

	return nil, exception.NewApiException(404, "user not found")
}

func (r *UserMemoryRepository) FindAndCheckJWT(claims *dto.DTOClaimsJwt) (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}

func (r *UserMemoryRepository) Insert(dtoRegisterRequest *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}

func (r *UserMemoryRepository) Update(dtoUserUpdate *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}

func (r *UserMemoryRepository) Delete(dtoUserDelete *dto.DTOUser) (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}
