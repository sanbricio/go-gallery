package userDTO

import (
	userEntity "go-gallery/src/domain/entities/user"
)

// UserDTO Representa los datos del usuario
// @Description Estructura que define los datos del usuario
type UserDTO struct {
	// Nombre de usuario
	// example "usuario123"
	Username string `json:"username" bson:"username" example:"usuario123"`

	// Contraseña
	// example "MiContraseñaSegura."
	Password string `json:"password" bson:"password" example:"MiContraseñaSegura."`

	// Correo electrónico
	// example "usuario@example.com"
	Email string `json:"email" bson:"email" example:"usuario@example.com"`

	// Apellido
	// example "Pérez"
	Lastname string `json:"lastname" bson:"lastname" example:"Pérez"`

	// Nombre
	// example "Juan"
	Firstname string `json:"firstname" bson:"firstname" example:"Juan"`
}

func FromUser(user *userEntity.User) *UserDTO {
	return &UserDTO{
		Username:  user.GetUsername(),
		Password:  user.GetPassword(),
		Email:     user.GetEmail(),
		Lastname:  user.GetLastname(),
		Firstname: user.GetFirstname(),
	}
}
