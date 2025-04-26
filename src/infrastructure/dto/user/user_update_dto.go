package userDTO

// UserUpdateDTO representa la estructura para actualizar los datos del usuario
// @Description Datos que pueden ser actualizados del usuario existente
type UserUpdateDTO struct {
	// Contraseña
	// example "NuevaContraseñaSegura."
	Password string `json:"password" example:"NuevaContraseñaSegura."`

	// Correo electrónico
	// example "nuevo.email@example.com"
	Email string `json:"email" example:"nuevo.email@example.com"`

	// Apellido
	// example "Gómez"
	Lastname string `json:"lastname" example:"Gómez"`

	// Nombre
	// example "Carlos"
	Firstname string `json:"firstname" example:"Carlos"`
}
