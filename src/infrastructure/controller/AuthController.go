package controller

import (
	"fmt"
	"go-gallery/src/commons/exception"
	"go-gallery/src/infrastructure/controller/handler"
	"go-gallery/src/infrastructure/controller/middlewares"
	"go-gallery/src/infrastructure/dto"
	"go-gallery/src/service"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userService        *service.UserService
	emailSenderService *service.EmailSenderService
	jwtMiddleware      *middlewares.JWTMiddleware
}

func NewAuthController(userService *service.UserService, emailSenderService *service.EmailSenderService, jwtMiddleware *middlewares.JWTMiddleware) *AuthController {
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

// @Summary      Iniciar sesión
// @Description  Autentica un usuario y genera un token JWT para guardarlo en una cookie
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.DTOLoginRequest true "Datos de autenticación"
// @Success      200 {object} dto.DTOLoginResponse "Se ha iniciado sesion correctamente"
// @Header       200 {string} Set-Cookie "Authorization=auth_token; HttpOnly; Secure"
// @Failure      400 {object} exception.ApiException "Contraseña incorrecta"
// @Failure      401 {object} exception.ApiException "No autorizado"
// @Failure      404 {object} exception.ApiException "Usuario no encontrado"
// @Failure      500 {object} exception.ApiException "Ha ocurrido un error inesperado"
// @Router       /auth/login [post]
func (c *AuthController) login(ctx *fiber.Ctx) error {
	dtoLoginRequest := new(dto.DTOLoginRequest)
	err := ctx.BodyParser(dtoLoginRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	dtoUser, errFind := c.userService.Find(dtoLoginRequest)
	if errFind != nil {
		return ctx.Status(errFind.Status).JSON(errFind)
	}

	errJWT := c.jwtMiddleware.CreateJWTToken(ctx, dtoUser.Username, dtoUser.Email)
	if errJWT != nil {
		return ctx.Status(errJWT.Status).JSON(errJWT.Status)
	}

	response := dto.DTOLoginResponse{
		Message:  "Se ha iniciado sesión correctamente",
		Username: dtoUser.Username,
		Email:    dtoUser.Email,
		Name:     dtoUser.Firstname,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

// @Summary      Registro de un nuevo usuario
// @Description  Registra un nuevo usuario en el sistema
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body dto.DTOUser true "Datos de registro"
// @Success      201 {object} dto.DTORegisterResponse "Usuario creado"
// @Failure      400 {object} exception.ApiException "Solicitud incorrecta"
// @Failure      500 {object} exception.ApiException "Ha ocurrido un error inesperado"
// @Router       /auth/register [post]
func (c *AuthController) register(ctx *fiber.Ctx) error {
	dtoRegisterRequest := new(dto.DTOUser)
	err := ctx.BodyParser(dtoRegisterRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	errHandler := handler.ProcessUser(dtoRegisterRequest.Password, dtoRegisterRequest.Email)
	if errHandler != nil {
		return ctx.Status(errHandler.Status).JSON(errHandler)
	}

	dtoUser, errInsert := c.userService.Insert(dtoRegisterRequest)
	if errInsert != nil {
		return ctx.Status(errInsert.Status).JSON(errInsert)
	}

	dto := dto.DTORegisterResponse{
		Username:  dtoUser.Username,
		Firstname: dtoUser.Firstname,
		Message:   "Se ha creado el usuario correctamente",
	}
	return ctx.Status(fiber.StatusCreated).JSON(dto)
}

// @Summary      Cerrar sesión
// @Description  Cierra la sesión del usuario autenticado, elimina la cookie auth_token
// @Tags         auth
// @Security     CookieAuth
// @Success      200 {object} dto.DTOMessageResponse "Se ha cerrado sesión correctamente"
// @Failure      401 {object} exception.ApiException "Usuario no autenticado"
// @Failure      403 {object} exception.ApiException "Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure      404 {object} exception.ApiException "Usuario no encontrado"
// @Failure      500 {object} exception.ApiException "Ha ocurrido un error inesperado"
// @Router       /auth/logout [post]
func (c *AuthController) logout(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
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
	return ctx.Status(fiber.StatusOK).JSON(&dto.DTOMessageResponse{
		Message: fmt.Sprintf("Se ha cerrado sesión correctamente, %s", claims.Username),
	})
}

// @Summary      Actualizar usuario
// @Description  Actualiza los datos de un usuario autenticado
// @Tags         auth
// @Security     CookieAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.DTOUpdateUser true "Datos de actualización"
// @Success      200 {object} dto.DTOMessageResponse "Se han actualizado los datos del usuario correctamente."
// @Failure      400 {object} exception.ApiException "Solicitud incorrecta"
// @Failure      401 {object} exception.ApiException "Usuario no autenticado"
// @Failure      403 {object} exception.ApiException "Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure      404 {object} exception.ApiException "Usuario no encontrado"
// @Failure      500 {object} exception.ApiException "Ha ocurrido un error inesperado"
// @Router       /auth/update [put]
func (c *AuthController) update(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtenemos los datos que quiere cambiar el usuario
	dtoUserUpdate := new(dto.DTOUpdateUser)
	err := ctx.BodyParser(dtoUserUpdate)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	// Creamos la entidad para proceder a la actualización del usuario
	dtoUser := &dto.DTOUser{
		Username:  claims.Username,
		Email:     dtoUserUpdate.Email,
		Password:  dtoUserUpdate.Password,
		Lastname:  dtoUserUpdate.Lastname,
		Firstname: dtoUserUpdate.Firstname,
	}

	// Si el usuario quiere cambiar el email y la contraseña lo validamos
	errUser = handler.ProcessUser(dtoUser.Password, dtoUser.Email)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Flag que nos permitira saber si necesitamos crear otro token JWT
	emailChanged := dtoUserUpdate.Email != "" && dtoUserUpdate.Email != claims.Email

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
	return ctx.Status(fiber.StatusOK).JSON(&dto.DTOMessageResponse{
		Message: fmt.Sprintf("Se han actualizado los datos del usuario %s correctamente.", dtoUser.Username),
	})
}

// @Summary      Solicitar eliminación de cuenta
// @Description  Envía un código de verificación al correo para eliminar la cuenta
// @Tags         auth
// @Security     CookieAuth
// @Success      200 {object} dto.DTOMessageResponse "Se ha enviado un código de confirmación al correo electrónico"
// @Failure      401 {object} exception.ApiException "Usuario no autenticado"
// @Failure      403 {object} exception.ApiException "Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure      404 {object} exception.ApiException "Usuario no encontrado"
// @Failure      500 {object} exception.ApiException "Ha ocurrido un error inesperado"
// @Router       /auth/request-delete [post]
func (c *AuthController) requestDelete(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Generamos el código único temporal
	code := handler.GenerateCode(claims.Username)

	// Enviamos al usuario un correo electrónico con el código
	errEmail := c.emailSenderService.SendEmail(code, claims.Email)
	if errEmail != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&dto.DTOMessageResponse{
			Message: "No se ha podido enviar un código de confirmación al correo electrónico debido a un fallo interno del sistema.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&dto.DTOMessageResponse{
		Message: fmt.Sprintf("Se ha enviado un código de confirmación al correo electrónico %s .", claims.Email),
	})
}

// @Summary      Confirmar eliminación de cuenta
// @Description  Elimina la cuenta de usuario tras verificar el código enviado
// @Tags         auth
// @Security     ApiKeyAuth
// @Accept       json
// @Produce      json
// @Param        request body dto.DTODeleteUser true "Datos para confirmar eliminación"
// @Success      200 {object} dto.DTOMessageResponse "Se han eliminado los datos del usuario correctamente"
// @Failure      400 {object} exception.ApiException "Solicitud incorrecta"
// @Failure      401 {object} exception.ApiException "Usuario no autenticado"
// @Failure      403 {object} exception.ApiException "Los datos proporcionados no coinciden con el usuario autenticado"
// @Failure      404 {object} exception.ApiException "Usuario no encontrado"
// @Failure      500 {object} exception.ApiException "Ha ocurrido un error inesperado"
// @Router       /auth/delete [delete]
func (c *AuthController) confirmDelete(ctx *fiber.Ctx) error {
	claims, ok := ctx.Locals("user").(*dto.DTOClaimsJwt)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Usuario no autenticado"))
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtenemos los datos de el usuario
	dtoDeleteUser := new(dto.DTODeleteUser)
	err := ctx.BodyParser(dtoDeleteUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	// Comprobamos que el código sea correcto
	ok = handler.VerifyCode(claims.Username, dtoDeleteUser.Code)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "El código de verificación es incorrecto"))
	}

	// Creamos la entidad para proceder a la eliminacion del usuario
	dtoUser := &dto.DTOUser{
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
	handler.RemoveCode(dtoUser.Username)

	// Eliminamos la cookie de autenticación
	c.jwtMiddleware.DeleteAuthCookie(ctx)

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(&dto.DTOMessageResponse{
		Message: fmt.Sprintf("Se han eliminado los datos del usuario %s correctamente.", dtoUser.Username),
	})
}
