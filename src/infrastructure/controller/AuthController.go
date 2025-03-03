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
	userService    *service.UserService
	authMiddleware *middlewares.AuthMiddleware
}

func NewAuthController(userService *service.UserService, authMiddleware *middlewares.AuthMiddleware) *AuthController {
	return &AuthController{
		userService:    userService,
		authMiddleware: authMiddleware,
	}
}

func (c *AuthController) SetUpRoutes(router fiber.Router) {
	router.Post("/login", c.login)
	router.Post("/register", c.register)
	router.Post("/logout", c.logout)
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

	token, errJWT := c.authMiddleware.CreateJwtToken(dtoUser.Username, dtoUser.Email, dtoUser.Firstname)
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
	cookie := ctx.Cookies("auth_token")
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "No se ha encontrado una sesión activa"))
	}

	// Obtener claims del token
	dto, err := c.authMiddleware.GetJWTClaimsFromCookie(cookie)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(exception.NewApiException(fiber.StatusUnauthorized, "Token inválido"))
	}
 
	// Eliminamos la cookie
	c.authMiddleware.DeleteAuthCookie(ctx)

	// Respuesta con el nombre del usuario
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Se ha cerrado sesión correctamente, %s", dto.Username),
	})
}
