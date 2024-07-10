package repository

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
)

type RepositoryUserMemory struct {
}

func Find(email, password string) (*entity.User, *exception.ApiException) {
	// Here you would query your database for the user with the given email or username
	if email == "test@mail.com" && password == "test12345" {
		return entity.NewUser("test", "test12345", "test@mail.com", "Prueba", "Prueba2"), nil
	}
	// hacer el documento y que lea del archivo que cree yo 
	return nil, exception.NewApiException(404, "user not found")
}
