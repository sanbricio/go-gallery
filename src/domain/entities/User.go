package entity

import (
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	username  string
	password  string
	email     string
	lastname  string
	firstname string
}

func NewUser(username, password, email, lastname, firstname string) *User {
	user := &User{
		username:  username,
		email:     email,
		lastname:  lastname,
		firstname: firstname,
	}

	user.hashPassword(password)

	return user
}

func NewUserFromDTO(username, password, email, lastname, firstname string) *User {
	if !isHashed(password) {
		log.Println("La contraseña del DTO no esta hasheada.")
		log.Println("Realizando hash de la contraseña")
		return NewUser(username, password, email, lastname, firstname)
	}

	user := &User{
		username:  username,
		password:  password,
		email:     email,
		lastname:  lastname,
		firstname: firstname,
	}

	return user
}

func (u *User) hashPassword(password string) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error al generar el hash de la contraseña: %v", err)
		return
	}
	u.password = string(hashedPassword)
}

func isHashed(password string) bool {
	return strings.HasPrefix(password, "$2a$") || strings.HasPrefix(password, "$2b$") || len(password) == 60
}

func (u *User) CheckPassword(password string) (*User, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (u *User) GetUsername() string {
	return u.username
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetLastname() string {
	return u.lastname
}

func (u *User) GetFirstname() string {
	return u.firstname
}
