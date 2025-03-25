package dto

type DTOUpdateUser struct {
	Password  string `json:"password"`
	Email     string `json:"email"`
	Lastname  string `json:"lastname"`
	Firstname string `json:"firstname"`
}
