package controller

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/controller/handler"
	"api-upload-photos/src/infrastructure/dto"
	"api-upload-photos/src/service"
	"time"

	"github.com/gofiber/fiber/v2"

	jtoken "github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	userService *service.UserService
	jwtSecret   string
}

func NewAuthController(userService *service.UserService, jwtSecret string) *AuthController {
	return &AuthController{
		userService: userService,
		jwtSecret:   jwtSecret,
	}
}

func (c *AuthController) SetUpRoutes(router fiber.Router) {
	router.Post("/login", c.loginUser)
	router.Post("/register", c.registerUser)
}

func (c *AuthController) loginUser(ctx *fiber.Ctx) error {
	dtoLoginRequest := new(dto.DTOUser)
	err := ctx.BodyParser(dtoLoginRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erróneo"))
	}

	dtoUser, errFind := c.userService.Find(dtoLoginRequest)
	if errFind != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(errFind)
	}

	day := time.Hour * 24

	// Crear las reclamaciones del JWT, incluyendo el usuario y la expiración
	claims := jtoken.MapClaims{
		"username": dtoUser.Username,
		"email":    dtoUser.Email,
		"name":     dtoUser.Firstname,
		"exp":      time.Now().Add(day * 1).Unix(),
	}

	// Crear el token
	token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

	// Firmar el token y devolverlo como respuesta
	t, err := token.SignedString([]byte(c.jwtSecret))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(exception.NewApiException(fiber.StatusInternalServerError, "Error al generar el token"))
	}

	// Devolver el token
	return ctx.Status(fiber.StatusOK).JSON(dto.DTOLoginResponse{
		Token: t,
	})
}

func (c *AuthController) registerUser(ctx *fiber.Ctx) error {
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
