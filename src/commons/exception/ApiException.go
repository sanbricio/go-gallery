package exception

type ApiException struct {
	status  int
	message string
}

func NewApiException(status int, message string) *ApiException {
	return &ApiException{
		status:  status,
		message: message,
	}
}

func (e *ApiException) GetStatus() int {
	return e.status
}

func (e *ApiException) GetMessage() string {
	return e.message
}
