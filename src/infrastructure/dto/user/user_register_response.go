package userDTO

// UserRegisterResponseDTO representa la respuesta tras registrar un usuario exitosamente
// @Description Respuesta generada después de crear un nuevo usuario
type UserRegisterResponseDTO struct {
	// Mensaje de confirmación
	// example "Se ha creado el usuario correctamente"
	Message string `json:"message" example:"Se ha creado el usuario correctamente"`

	// Nombre de usuario
	// example "usuario123"
	Username string `json:"username" example:"usuario123"`

	// Nombre
	// example "Juan"
	Firstname string `json:"firstname" example:"Juan"`
}
