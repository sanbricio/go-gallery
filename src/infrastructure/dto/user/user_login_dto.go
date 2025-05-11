package userDTO

// LoginRequestDTO representa la estructura para solicitar inicio de sesión
// @Description Datos requeridos para realizar autentificación del usuario
type LoginRequestDTO struct {
	// Nombre de usuario
	// example "usuario123"
	Username string `json:"username"  bson:"username" example:"usuario123"`

	// Contraseña del usuario
	// example "MiContraseñaSegura."
	Password string `json:"password" bson:"pasword" example:"MiContraseñaSegura."`
}

// LoginResponseDTO representa la respuesta de autenticación exitosa
// @Description Respuesta al iniciar sesión correctamente
type LoginResponseDTO struct {
	// Mensaje de confirmación
	// example "Se ha iniciado sesión correctamente"
	Message string `json:"message" example:"Se ha iniciado sesión correctamente"`

	// Nombre de usuario autenticado
	// example "usuario123"
	Username string `json:"username" example:"usuario123"`

	// Correo electrónico del usuario
	// example "usuario@example.com"
	Email string `json:"email" example:"usuario@example.com"`

	// Nombre del usuario
	// example "Juan"
	Firstname string `json:"firstname" example:"Juan"`

	// Apellido del usuario
	// example "Bricio"
	Lastname string `json:"lastname" example:"Bricio"`
}
