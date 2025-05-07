package dto

// MessageResponseDTO representa una respuesta con un mensaje genérico
// @Description Respuesta con un mensaje para informar al usuario que ha ocurrido
type MessageResponseDTO struct {
	// Mensaje de respuesta
	// example "Operación realizada con éxito"
	Message string `json:"message" example:"Ha funcionado correctamente"`
}
