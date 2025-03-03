package dto

type DTOClaimsJwt struct {
	Username   string `json:"username"`
	Email      string `json:"email"`
	Firstname  string `json:"firstname"`
	Expiration int64  `json:"expiration"`
}
