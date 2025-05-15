package userDTO

// UserDeleteDTO representa los datos requeridos para eliminar un usuario
// @Description Datos necesarios para proceder con la eliminación del usuario
type UserDeleteDTO struct {
	// Contraseña del usuario para confirmar la eliminación
	Password string `json:"password" example:"MiContraseñaSegura"`

	// Código de verificación para la eliminación
	Code string `json:"code" example:"123456"`
}
