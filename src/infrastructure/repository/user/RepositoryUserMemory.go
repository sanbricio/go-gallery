package repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"api-upload-photos/src/infrastructure/dto"
	"os"
)

type RepositoryUserMemory struct {
}

func (r *RepositoryUserMemory) Find(username, password string) (*dto.DTOUser, *exception.ApiException) {

	files, err := os.ReadDir("data")
	if err != nil {
		return nil, exception.NewApiException(500, "Error al leer el directorio")
	}

	for _, file := range files {
		if file.Name() == "prueba.json" {
			if username == "test" && password == "test12345" {
				entity := entity.NewUser("test", "test12345", "test@mail.com", "Prueba", "Prueba2")
				dto := dto.FromUser(entity)
				return dto,nil
			}
		}
	}

	return nil, exception.NewApiException(404, "user not found")
}

func (r *RepositoryUserMemory) Insert() (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}

func (r *RepositoryUserMemory) Update() (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}

func (r *RepositoryUserMemory) Delete() (*dto.DTOUser, *exception.ApiException) {
	return nil, nil
}
