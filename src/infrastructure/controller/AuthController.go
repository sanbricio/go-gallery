package controller

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/controller/handler"
	"api-upload-photos/src/infrastructure/controller/middlewares"
	"api-upload-photos/src/infrastructure/dto"
	"api-upload-photos/src/service"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	userService        *service.UserService
	emailSenderService *service.EmailSenderService
	authMiddleware     *middlewares.AuthMiddleware
}

func NewAuthController(userService *service.UserService, emailSenderService *service.EmailSenderService, authMiddleware *middlewares.AuthMiddleware) *AuthController {
	return &AuthController{
		userService:        userService,
		emailSenderService: emailSenderService,
		authMiddleware:     authMiddleware,
	}
}

func (c *AuthController) SetUpRoutes(router fiber.Router) {
	router.Post("/login", c.login)
	router.Post("/register", c.register)
	router.Post("/logout", c.logout)
	router.Put("/update", c.update)
	router.Post("/request-delete", c.requestDelete)
	router.Delete("/delete", c.confirmDelete)
}

func (c *AuthController) login(ctx *fiber.Ctx) error {
	dtoLoginRequest := new(dto.DTOUser)
	err := ctx.BodyParser(dtoLoginRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	dtoUser, errFind := c.userService.Find(dtoLoginRequest)
	if errFind != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errFind)
	}

	token, errJWT := c.authMiddleware.CreateJwtToken(dtoUser.Username, dtoUser.Email)
	if errJWT != nil {
		return ctx.Status(errJWT.Status).JSON(errJWT)
	}

	// Configuramos la cookie con el JWT
	ctx.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(2 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	response := dto.DTOLoginResponse{
		Message:  "Se ha iniciado sesión correctamente",
		Username: dtoUser.Username,
		Email:    dtoUser.Email,
		Name:     dtoUser.Firstname,
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *AuthController) register(ctx *fiber.Ctx) error {
	dtoRegisterRequest := new(dto.DTOUser)
	err := ctx.BodyParser(dtoRegisterRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	errHandler := handler.ProcessUser(dtoRegisterRequest)
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

func (c *AuthController) logout(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(c.authMiddleware.GetCookieName())
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No se ha encontrado una sesión activa"))
	}

	// Obtener claims del token
	claims, err := c.authMiddleware.GetJWTClaimsFromCookie(cookie)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Token inválido"))
	}

	// Eliminamos la cookie
	c.authMiddleware.DeleteAuthCookie(ctx)

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Se ha cerrado sesión correctamente, %s", claims.Username),
	})
}

func (c *AuthController) update(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(c.authMiddleware.GetCookieName())
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No se ha encontrado una sesión activa"))
	}

	// Obtener claims del token
	claims, errJWT := c.authMiddleware.GetJWTClaimsFromCookie(cookie)
	if errJWT != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errJWT)
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtenemos los datos que quiere cambiar el usuario
	dtoUserUpdate := new(dto.DTOUser)
	err := ctx.BodyParser(dtoUserUpdate)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	// Si el usuario quiere cambiar la contraseña comprobamos que esta no este vacia y la validamos
	if dtoUserUpdate.Password != "" {
		errHandler := handler.ValidatePassword(dtoUserUpdate.Password)
		if errHandler != nil {
			return ctx.Status(errHandler.Status).JSON(errHandler)
		}
	}

	// Actualizamos el usuario
	_, errUpdate := c.userService.Update(dtoUserUpdate)
	if errUpdate != nil {
		return ctx.Status(errUpdate.Status).JSON(errUpdate)
	}

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Se han actualizado los datos del usuario %s correctamente.", dtoUserUpdate.Username),
	})
}

func (c *AuthController) requestDelete(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies(c.authMiddleware.GetCookieName())
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No se ha encontrado una sesión activa"))
	}

	// Obtenemos claims del token
	claims, errJWT := c.authMiddleware.GetJWTClaimsFromCookie(cookie)
	if errJWT != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errJWT)
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Generamos el código único temporal
	code := handler.GenerateCode(claims.Email)

	// Enviamos al usuario un correo electrónico con el código
	errEmail := c.emailSenderService.SendEmail(code, claims.Email)
	if errEmail != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "No se ha podido enviar un código de confirmación al correo electrónico debido a un fallo interno del sistema.",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Se ha enviado un código de confirmación al correo electrónico %s .", claims.Email),
	})
}

func (c *AuthController) confirmDelete(ctx *fiber.Ctx) error {
	// TODO Aqui faltaria comprobar el código de confirmacion para la eliminacion de la cuenta
	cookie := ctx.Cookies(c.authMiddleware.GetCookieName())
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No se ha encontrado una sesión activa"))
	}

	// Obtener claims del token
	claims, errJWT := c.authMiddleware.GetJWTClaimsFromCookie(cookie)
	if errJWT != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errJWT)
	}

	// Comprobamos que el usuario existe
	_, errUser := c.userService.FindAndCheckJWT(claims)
	if errUser != nil {
		return ctx.Status(errUser.Status).JSON(errUser)
	}

	// Obtenemos los datos de el usuario
	dtoUser := new(dto.DTOUser)
	err := ctx.BodyParser(dtoUser)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	// Eliminamos el usuario
	_, errDelete := c.userService.Delete(dtoUser)
	if errDelete != nil {
		return ctx.Status(errDelete.Status).JSON(errDelete)
	}

	// Eliminamos la cookie de autenticación
	c.authMiddleware.DeleteAuthCookie(ctx)

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Se han eliminado los datos del usuario %s correctamente.", dtoUser.Username),
	})
}
