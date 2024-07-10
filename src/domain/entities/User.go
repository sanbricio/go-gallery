package entity

type User struct {
	username  string
	password  string
	email     string
	lastname  string
	firstname string
}

func NewUser(username string, password string, email string, lastname string, firstname string) *User {
	return &User{
		username:  username,
		password:  password,
		email:     email,
		lastname:  lastname,
		firstname: firstname,
	}
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

func (u* User) GetLastname() string{
	return u.lastname
}

func (u* User) GetFirstname() string {
	return u.firstname
}
