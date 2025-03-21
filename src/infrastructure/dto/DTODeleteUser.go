package dto

type DTODeleteUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	Code     string `json:"code"`
}
