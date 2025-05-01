package userController

import (
	"fmt"
	"go-gallery/src/commons/exception"

	codeGeneratorHandler "go-gallery/src/infrastructure/controller/handler/codeGenerator"
	userHandler "go-gallery/src/infrastructure/controller/user/handler"
	userMiddleware "go-gallery/src/infrastructure/controller/user/middlewares"
	userDTO "go-gallery/src/infrastructure/dto/user"
	emailService "go-gallery/src/service/email"
	userService "go-gallery/src/service/user"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userService        *userService.UserService
	emailSenderService *emailService.EmailSenderService
	jwtMiddleware      *userMiddleware.JWTMiddleware
}

func NewAuthController(userService *userService.UserService, emailSenderService *emailService.EmailSenderService, jwtMiddleware *userMiddleware.JWTMiddleware) *AuthController {
	return &AuthController{
		userService:        userService,
		emailSenderService: emailSenderService,
		jwtMiddleware:      jwtMiddleware,
	}
}

func (c *AuthController) SetUpRoutes(router fiber.Router) {
	router.Post("/login", c.login)
	router.Post("/register", c.register)
	router.Post("/logout", c.jwtMiddleware.Handler(), c.logout)
	router.Put("/update", c.jwtMiddleware.Handler(), c.update)
	router.Post("/request-delete", c.jwtMiddleware.Handler(), c.requestDelete)
	router.Delete("/delete", c.jwtMiddleware.Handler(), c.confirmDelete)
}

// @Summary		Iniciar sesión
// @Description	Autentica un usuario y genera un token JWT para guardarlo en una cookie
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		userDTO.LoginRequestDTO		true	"Datos de autenticación"
// @Success		200		{object}	userDTO.LoginResponseDTO	"Se ha iniciado sesion correctamente"
// @Header			200		{string}	Set-Cookie					"Authorization=auth_token; HttpOnly; Secure"
// @Failure		400		{object}	exception.ApiException		"Contraseña incorrecta"
// @Failure		401		{object}	exception.ApiException		"No autorizado"
// @Failure		404		{object}	exception.ApiException		"Usuario no encontrado"
// @Failure		500		{object}	exception.ApiException		"Ha ocurrido un error inesperado"
// @Router			/auth/login [post]
func (c *AuthController) login(ctx *fiber.Ctx) error {
	loginRequestDTO := new(userDTO.LoginRequestDTO)
	err := ctx.BodyParser(loginRequestDTO)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	user, errFind := c.userService.Find(loginRequestDTO)
	if errFind != nil {
		return ctx.Status(errFind.Status).JSON(errFind)
	}

	errJWT := c.jwtMiddleware.CreateJWTToken(ctx, user.Username, user.Email)
	if errJWT != nil {
		return ctx.Status(errJWT.Status).JSON(errJWT.Status)
	}

	responseDTO := userDTO.LoginResponseDTO{
		Message:  "Se ha iniciado sesión correctamente",
		Username: user.Username,
		Email:    user.Email,
		Name:     user.Firstname,
	}

	return ctx.Status(fiber.StatusOK).JSON(responseDTO)
}

// @Summary		Registro de un nuevo usuario
// @Description	Registra un nuevo usuario en el sistema
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			request	body		userDTO.UserDTO					true	"Datos de registro"
// @Success		201		{object}	userDTO.UserRegisterResponseDTO	"Usuario creado"
// @Failure		400		{object}	exception.ApiException			"Solicitud incorrecta"
// @Failure		500		{object}	exception.ApiException			"Ha ocurrido un error inesperado"
// @Router			/auth/register [post]
func (c *AuthController) register(ctx *fiber.Ctx) error {
	registerRequestDTO := new(userDTO.UserDTO)
	err := ctx.BodyParser(registerRequestDTO)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	errHandler := userHandler.ProcessUser(registerRequestDTO.Password, registerRequestDTO.Email)
	if errHandler != nil {
		return ctx.Status(errHandler.Status).JSON(errHandler)
	}

	user, errInsert := c.userService.Insert(registerRequestDTO)
	if errInsert != nil {
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	dto := userDTO.UserRegisterResponseDTO{
		Username:  user.Username,
		Firstname: user.Firstname,
		Message:   "Se ha creado el usuario correctamente",
	}
	return ctx.Status(fiber.StatusCreated).JSON(dto)
}

// @Summary		Cerrar sesión
// @Description	Cierra la sesión del usuario autenticado, elimina la cookie auth_token
// @Tags			auth
// @Security		CookieAuth
// @Success		200	{object}	userDTO.UserRegisterResponseDTO	"Se ha cerrado sesión correctamente"
// @Failure		401	{object}	exception.ApiException			"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException			"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException			"Usuario no encontrado"
// @Failure		500	{object}	exception.ApiException			"Ha ocurrido un error inesperado"
// @Router			/auth/logout [post]
func (c *AuthController) logout(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}
	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Eliminamos la cookie
	c.jwtMiddleware.DeleteAuthCookie(ctx)

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserRegisterResponseDTO{
		Message: fmt.Sprintf("Se ha cerrado sesión correctamente, %s", claims.Username),
	})
}

// @Summary		Actualizar usuario
// @Description	Actualiza los datos de un usuario autenticado
// @Tags			auth
// @Security		CookieAuth
// @Accept			json
// @Produce		json
// @Param			request	body		userDTO.UserUpdateDTO			true	"Datos de actualización"
// @Success		200		{object}	userDTO.UserRegisterResponseDTO	"Se han actualizado los datos del usuario correctamente."
// @Failure		400		{object}	exception.ApiException			"Solicitud incorrecta"
// @Failure		401		{object}	exception.ApiException			"Usuario no autenticado"
// @Failure		403		{object}	exception.ApiException			"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404		{object}	exception.ApiException			"Usuario no encontrado"
// @Failure		500		{object}	exception.ApiException			"Ha ocurrido un error inesperado"
// @Router			/auth/update [put]
func (c *AuthController) update(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtenemos los datos que quiere cambiar el usuario
	user := new(userDTO.UserUpdateDTO)
	err := ctx.BodyParser(user)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	// Creamos la entidad para proceder a la actualización del usuario
	dtoUser := &userDTO.UserDTO{
		Username:  claims.Username,
		Email:     user.Email,
		Password:  user.Password,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
	}

	// Si el usuario quiere cambiar el email y la contraseña lo validamos
	errUser = userHandler.ProcessUser(dtoUser.Password, dtoUser.Email)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Flag que nos permitira saber si necesitamos crear otro token JWT
	emailChanged := user.Email != "" && user.Email != claims.Email

	// Actualizamos el usuario
	_, errUpdate := c.userService.Update(dtoUser)
	if errUpdate != nil {
		return ctx.Status(errUpdate.Status).JSON(errUpdate)
	}

	// Si el email ha cambiado, generamos un nuevo token JWT (debido a que el email es parte de la información del token)
	if emailChanged {
		errJWT := c.jwtMiddleware.CreateJWTToken(ctx, dtoUser.Username, dtoUser.Email)
		if errJWT != nil {
			return ctx.Status(errJWT.Status).JSON(errJWT.Status)
		}
	}

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserMessageResponseDTO{
		Message: fmt.Sprintf("Se han actualizado los datos del usuario %s correctamente.", dtoUser.Username),
	})
}

// @Summary		Solicitar eliminación de cuenta
// @Description	Envía un código de verificación al correo para eliminar la cuenta
// @Tags			auth
// @Security		CookieAuth
// @Success		200	{object}	userDTO.UserRegisterResponseDTO	"Se ha enviado un código de confirmación al correo electrónico"
// @Failure		401	{object}	exception.ApiException			"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException			"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException			"Usuario no encontrado"
// @Failure		500	{object}	exception.ApiException			"Ha ocurrido un error inesperado"
// @Router			/auth/request-delete [post]
func (c *AuthController) requestDelete(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Generamos el código único temporal
	code := codeGeneratorHandler.GenerateCode(claims.Username)

	// Enviamos al usuario un correo electrónico con el código
	errEmail := c.emailSenderService.SendEmail(code, claims.Email)
	if errEmail != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&userDTO.UserMessageResponseDTO{
			Message: "No se ha podido enviar un código de confirmación al correo electrónico debido a un fallo interno del sistema.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserMessageResponseDTO{
		Message: fmt.Sprintf("Se ha enviado un código de confirmación al correo electrónico %s .", claims.Email),
	})
}

// @Summary		Confirmar eliminación de cuenta
// @Description	Elimina la cuenta de usuario tras verificar el código enviado
// @Tags			auth
// @Security		ApiKeyAuth
// @Accept			json
// @Produce		json
// @Param			request	body		userDTO.UserDeleteDTO			true	"Datos para confirmar eliminación"
// @Success		200		{object}	userDTO.UserRegisterResponseDTO	"Se han eliminado los datos del usuario correctamente"
// @Failure		400		{object}	exception.ApiException			"Solicitud incorrecta"
// @Failure		401		{object}	exception.ApiException			"Usuario no autenticado"
// @Failure		403		{object}	exception.ApiException			"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404		{object}	exception.ApiException			"Usuario no encontrado"
// @Failure		500		{object}	exception.ApiException			"Ha ocurrido un error inesperado"
// @Router			/auth/delete [delete]
func (c *AuthController) confirmDelete(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtenemos los datos de el usuario
	dtoDeleteUser := new(userDTO.UserDeleteDTO)
	err := ctx.BodyParser(dtoDeleteUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	// Comprobamos que el código sea correcto
	ok = codeGeneratorHandler.VerifyCode(claims.Username, dtoDeleteUser.Code)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "El código de verificación es incorrecto"))
	}

	// Creamos la entidad para proceder a la eliminacion del usuario
	dtoUser := &userDTO.UserDTO{
		Username: claims.Username,
		Email:    claims.Email,
		Password: dtoDeleteUser.Password,
	}

	// Eliminamos el usuario
	_, errDelete := c.userService.Delete(dtoUser)
	if errDelete != nil {
		return ctx.Status(errDelete.Status).JSON(errDelete)
	}

	// Eliminamos el codigo unico del usuario del mapa
	codeGeneratorHandler.RemoveCode(dtoUser.Username)

	// Eliminamos la cookie de autenticación
	c.jwtMiddleware.DeleteAuthCookie(ctx)

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserMessageResponseDTO{
		Message: fmt.Sprintf("Se han eliminado los datos del usuario %s correctamente.", dtoUser.Username),
	})
}
