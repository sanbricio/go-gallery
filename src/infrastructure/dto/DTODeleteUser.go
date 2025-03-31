package dto

// DTODeleteUser representa los datos requeridos para eliminar un usuario
// @Description Datos necesarios para proceder con la eliminación del usuario
type DTODeleteUser struct {
	// Contraseña del usuario para confirmar la eliminación
	// example "MiContraseñaSegura"
	Password string `json:"password" example:"MiContraseñaSegura"`

	// Código de verificación para la eliminación
	// example "123456"
	Code string `json:"code" example:"123456"`
}
