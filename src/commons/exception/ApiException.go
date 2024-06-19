package exception

type ApiException struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewApiException(status int, message string) *ApiException {
	return &ApiException{
		Status:  status,
		Message: message,
	}
}
