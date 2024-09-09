package entity

import "golang.org/x/crypto/bcrypt"

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
		password:  password,
		lastname:  lastname,
		firstname: firstname,
	}
	return user
}

func (u *User) CheckPasswordIntegrity(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.password), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
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
