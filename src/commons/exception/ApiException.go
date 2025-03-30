package exception

// ApiException representa un error genérico con un estado y un mensaje
// @Description Estructura para manejar excepciones con un código de estado y un mensaje de error
type ApiException struct {
	// Código de estado HTTP
	// example 400
	Status int `json:"status"  example:"400"`

	// Mensaje de error
	// example "Solicitud incorrecta"
	Message string `json:"message"  example:"Solicitud incorrecta"`
}

func NewApiException(status int, message string) *ApiException {
	return &ApiException{
		Status:  status,
		Message: message,
	}
}
