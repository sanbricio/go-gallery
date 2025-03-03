package dto

type DTOLoginResponse struct {
	Message  string `json:"message"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
}
