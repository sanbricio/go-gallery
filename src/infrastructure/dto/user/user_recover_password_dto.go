package userDTO

// PasswordRecoveryRequestDTO representa la estructura para solicitar la recuperación de la contraseña
// @Description Datos requeridos para realizar la recuperación de la contraseña
type PasswordRecoveryRequestDTO struct {
	// Correo electrónico del usuario
	Email string `json:"email" example:"usuario@example.com"`
}


// PasswordRecoveryConfirmDTO representa la estructura para solicitar la recuperación de la contraseña
// @Description Datos requeridos para realizar la recuperación de la contraseña
type PasswordRecoveryConfirmDTO struct {
	// Correo electrónico del usuario
	Email string `json:"email" example:"usuario@example.com"`
	// Código de verificación para la eliminación
	Code string `json:"code" example:"123456"`
	// Nueva contraseña del usuario
	NewPassword string `json:"newPassword" example:"NuevaContraseñaSegura."`
}
