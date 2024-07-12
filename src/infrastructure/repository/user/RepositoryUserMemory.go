package repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"os"
)

type RepositoryUserMemory struct {
}

func Find(email, password string) (*entity.User, *exception.ApiException) {

	files, err := os.ReadDir("data")
	if err != nil {
		return nil, exception.NewApiException(500, "Error al leer el directorio")
	}

	for _, file := range files {
		if file.Name() == "prueba.json" {
			if email == "test@mail.com" && password == "test12345" {
				return entity.NewUser("test", "test12345", "test@mail.com", "Prueba", "Prueba2"), nil
			}
		}
	}

	return nil, exception.NewApiException(404, "user not found")
}
