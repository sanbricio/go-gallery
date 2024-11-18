package builder

import (
	"api-upload-photos/src/commons/exception"
	entity "api-upload-photos/src/domain/entities"
	"api-upload-photos/src/infrastructure/dto"
	"strings"
)

type UserBuilder struct {
	username  string
	password  string
	email     string
	lastname  string
	firstname string
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{}
}

func (b *UserBuilder) FromDTO(dto *dto.DTOUser) *UserBuilder {
	b.username = dto.Username
	b.password = dto.Password
	b.email = dto.Email
	b.lastname = dto.Lastname
	b.firstname = dto.Firstname

	return b
}

func (b *UserBuilder) Build() (*entity.User, *exception.BuilderException) {
	err := b.validateUser()
	if err != nil {
		return nil, err
	}

	if !b.isHashed(b.password) {
		hashedPassword, err := entity.HashPassword(b.password)
		if err != nil {
			return nil, exception.NewBuilderException("password", "Error al hashear la contraseña")
		}
		b.password = hashedPassword
	}

	return entity.NewUser(b.username, b.password, b.email, b.lastname, b.firstname), nil
}

func (b *UserBuilder) validateUser() *exception.BuilderException {
	if b.username == "" {
		return exception.NewBuilderException("username", "El campo 'username' no debe estar vacio")
	}

	if b.password == "" {
		return exception.NewBuilderException("password", "El campo 'password' no debe estar vacio")
	}

	if b.email == "" {
		return exception.NewBuilderException("email", "El campo 'email' no debe estar vacio")
	}

	if b.firstname == "" {
		return exception.NewBuilderException("username", "El campo 'firstname' no debe estar vacio")
	}

	return nil
}

func (b *UserBuilder) isHashed(password string) bool {
	return strings.HasPrefix(password, "$2a$") || strings.HasPrefix(password, "$2b$") || len(password) == 60
}

func (b *UserBuilder) SetUsername(username string) *UserBuilder {
	b.username = username
	return b
}

func (b *UserBuilder) SetPassword(password string) *UserBuilder {
	b.password = password
	return b
}

func (b *UserBuilder) SetEmail(email string) *UserBuilder {
	b.email = email
	return b
}

func (b *UserBuilder) SetLastname(lastname string) *UserBuilder {
	b.lastname = lastname
	return b
}

func (b *UserBuilder) SetFirstname(firstname string) *UserBuilder {
	b.firstname = firstname
	return b
}
