package userController

import (
	"fmt"
	"go-gallery/src/commons/exception"

	userHandler "go-gallery/src/infrastructure/controller/user/handler"
	userMiddleware "go-gallery/src/infrastructure/controller/user/middlewares"
	imageDTO "go-gallery/src/infrastructure/dto/image"
	userDTO "go-gallery/src/infrastructure/dto/user"
	log "go-gallery/src/infrastructure/logger"
	codeGeneratorService "go-gallery/src/service/codeGenerator"
	emailService "go-gallery/src/service/email"
	imageService "go-gallery/src/service/image"
	userService "go-gallery/src/service/user"

	"github.com/gofiber/fiber/v2"
)

var logger log.Logger

const (
	INVALID_LOGIN_REQUEST_MSG    string = "Invalid JSON in the request body"
	INVALID_AUTHENTIFICATION_MSG string = "User not authenticated"
	CLAIMS_NOT_FOUND_MSG         string = "Unauthorized: no user claims found"
	PREFIX_DELETE_CODE_GENERATOR string = "delete"
)

type AuthController struct {
	userService          *userService.UserService
	emailSenderService   *emailService.EmailSenderService
	imageService         *imageService.ImageService
	codeGeneratorService *codeGeneratorService.CodeGeneratorService
	jwtMiddleware        *userMiddleware.JWTMiddleware
}

func NewAuthController(userService *userService.UserService, emailSenderService *emailService.EmailSenderService,
	imageService *imageService.ImageService, codeGeneratorService *codeGeneratorService.CodeGeneratorService,
	jwtMiddleware *userMiddleware.JWTMiddleware) *AuthController {
	logger = log.Instance()
	return &AuthController{
		userService:          userService,
		emailSenderService:   emailSenderService,
		imageService:         imageService,
		codeGeneratorService: codeGeneratorService,
		jwtMiddleware:        jwtMiddleware,
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
	logger.Info("POST /login called")

	loginRequestDTO := new(userDTO.LoginRequestDTO)
	err := ctx.BodyParser(loginRequestDTO)
	if err != nil {
		logger.Error("Invalid JSON in login request")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, INVALID_LOGIN_REQUEST_MSG))
	}

	user, errFind := c.userService.Find(loginRequestDTO)
	if errFind != nil {
		logger.Error(fmt.Sprintf("Error finding user: %s", errFind.Message))
		return ctx.Status(errFind.Status).JSON(errFind)
	}

	errJWT := c.jwtMiddleware.CreateJWTToken(ctx, user.Username, user.Email)
	if errJWT != nil {
		logger.Error(fmt.Sprintf("Error creating JWT token: %s", errJWT.Message))
		return ctx.Status(errJWT.Status).JSON(errJWT.Status)
	}

	responseDTO := userDTO.LoginResponseDTO{
		Message:   "Login successful",
		Username:  user.Username,
		Email:     user.Email,
		Firstname: user.Firstname,
		Lastname:  user.Lastname,
	}

	logger.Info(fmt.Sprintf("User %s logged in successfully", user.Username))
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
	logger.Info("POST /register called")

	registerRequestDTO := new(userDTO.UserDTO)
	err := ctx.BodyParser(registerRequestDTO)
	if err != nil {
		logger.Error("Invalid JSON in register request")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, INVALID_LOGIN_REQUEST_MSG))
	}

	errHandler := userHandler.ProcessUser(registerRequestDTO.Password, registerRequestDTO.Email)
	if errHandler != nil {
		logger.Error(fmt.Sprintf("Error processing user data: %s", errHandler.Message))
		return ctx.Status(errHandler.Status).JSON(errHandler)
	}

	user, errInsert := c.userService.Insert(registerRequestDTO)
	if errInsert != nil {
		logger.Error(fmt.Sprintf("Error inserting new user: %s", errInsert.Message))
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	dto := userDTO.UserRegisterResponseDTO{
		Username:  user.Username,
		Firstname: user.Firstname,
		Message:   "User registered successfully",
	}

	logger.Info(fmt.Sprintf("User %s registered successfully", user.Username))
	return ctx.Status(fiber.StatusCreated).JSON(dto)
}

// @Summary		Cerrar sesión
// @Description	Cierra la sesión del usuario autenticado, elimina la cookie auth_token
// @Tags			auth
// @Security		CookieAuth
// @Success		200	{object}	userDTO.UserMessageResponseDTO	"Se ha cerrado sesión correctamente"
// @Failure		401	{object}	exception.ApiException			"Usuario no autenticado"
// @Failure		403	{object}	exception.ApiException			"Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure		404	{object}	exception.ApiException			"Usuario no encontrado"
// @Failure		500	{object}	exception.ApiException			"Ha ocurrido un error inesperado"
// @Router			/auth/logout [post]
func (c *AuthController) logout(ctx *fiber.Ctx) error {
	logger.Info("POST /logout called")

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(CLAIMS_NOT_FOUND_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	c.jwtMiddleware.DeleteAuthCookie(ctx)
	logger.Info(fmt.Sprintf("User %s logged out successfully", claims.Username))

	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserMessageResponseDTO{
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
	logger.Info("PUT /update called")

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(CLAIMS_NOT_FOUND_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	user := new(userDTO.UserUpdateDTO)
	err := ctx.BodyParser(user)
	if err != nil {
		logger.Error("Invalid JSON in update request")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, INVALID_LOGIN_REQUEST_MSG))
	}

	dtoUser := &userDTO.UserDTO{
		Username:  claims.Username,
		Email:     user.Email,
		Password:  user.Password,
		Lastname:  user.Lastname,
		Firstname: user.Firstname,
	}

	errUser := userHandler.ProcessUser(dtoUser.Password, dtoUser.Email)
	if errUser != nil {
		logger.Error(fmt.Sprintf("Error processing user data: %s", errUser.Message))
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	emailChanged := user.Email != "" && user.Email != claims.Email

	_, errUpdate := c.userService.Update(dtoUser)
	if errUpdate != nil {
		logger.Error(fmt.Sprintf("Error updating user: %s", errUpdate.Message))
		return ctx.Status(errUpdate.Status).JSON(errUpdate)
	}

	if emailChanged {
		errJWT := c.jwtMiddleware.CreateJWTToken(ctx, dtoUser.Username, dtoUser.Email)
		if errJWT != nil {
			logger.Error(fmt.Sprintf("Error creating new JWT token: %s", errJWT.Message))
			return ctx.Status(errJWT.Status).JSON(errJWT.Status)
		}
		logger.Info(fmt.Sprintf("User %s email updated, new JWT token created", dtoUser.Username))
	}

	logger.Info(fmt.Sprintf("User %s updated successfully", dtoUser.Username))
	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserMessageResponseDTO{
		Message: fmt.Sprintf("User %s updated successfully.", dtoUser.Username),
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
	logger.Info("POST /request-delete called")

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(CLAIMS_NOT_FOUND_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	code, err := c.codeGeneratorService.GenerateCode("delete", claims.Username)
	if err != nil {
		logger.Error(fmt.Sprintf("Error generating delete code: %v", err))
		return ctx.Status(fiber.StatusInternalServerError).JSON(exception.NewApiException(fiber.StatusInternalServerError, "Error generating delete code"))
	}

	errEmail := c.emailSenderService.SendEmail(code, claims.Email)
	if errEmail != nil {
		logger.Error(fmt.Sprintf("Failed to send delete code email to %s", claims.Email))
		return ctx.Status(fiber.StatusInternalServerError).JSON(&userDTO.UserMessageResponseDTO{
			Message: "Failed to send confirmation code email due to an internal system error.",
		})
	}

	logger.Info(fmt.Sprintf("Delete code sent successfully to email: %s", claims.Email))
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
	logger.Info("DELETE /delete called")

	claims, ok := ctx.Locals("user").(*userDTO.JwtClaimsDTO)
	if !ok {
		logger.Error(CLAIMS_NOT_FOUND_MSG)
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, INVALID_AUTHENTIFICATION_MSG))
	}

	dtoDeleteUser := new(userDTO.UserDeleteDTO)
	err := ctx.BodyParser(dtoDeleteUser)
	if err != nil {
		logger.Error("Invalid JSON in delete request")
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, INVALID_LOGIN_REQUEST_MSG))
	}

	ok = c.codeGeneratorService.VerifyCode("delete", claims.Username, dtoDeleteUser.Code)
	if !ok {
		logger.Error("Invalid verification code")
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Invalid verification code"))
	}

	// Create request to delete all user images/thumbnails by owner
	dtoDeleteImage := new(imageDTO.ImageDeleteRequestDTO)
	dtoDeleteImage.Owner = claims.Username

	_, errImageResponse := c.imageService.DeleteAll(dtoDeleteImage)
	if errImageResponse != nil {
		logger.Error(fmt.Sprintf("Error deleting all images for user %s: %s", claims.Username, errImageResponse.Message))
	}

	logger.Info(fmt.Sprintf("All images/thumbnails for user %s deleted successfully", claims.Username))

	dtoUser := &userDTO.UserDTO{
		Username: claims.Username,
		Email:    claims.Email,
		Password: dtoDeleteUser.Password,
	}

	_, errDelete := c.userService.Delete(dtoUser)
	if errDelete != nil {
		logger.Error(fmt.Sprintf("Error deleting user: %s", errDelete.Message))
		return ctx.Status(errDelete.Status).JSON(errDelete)
	}

	c.jwtMiddleware.DeleteAuthCookie(ctx)

	logger.Info(fmt.Sprintf("User %s deleted successfully", dtoUser.Username))
	return ctx.Status(fiber.StatusOK).JSON(&userDTO.UserMessageResponseDTO{
		Message: fmt.Sprintf("User %s deleted successfully.", dtoUser.Username),
	})
}
