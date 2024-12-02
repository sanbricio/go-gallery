package controller

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/infrastructure/controller/handler"
	"api-upload-photos/src/infrastructure/controller/middlewares"
	"api-upload-photos/src/infrastructure/dto"
	"api-upload-photos/src/service"
	"time"

	"github.com/gofiber/fiber/v2"
	jtoken "github.com/golang-jwt/jwt/v5"
)

type Controller struct {
	app          *fiber.App
	serviceImage *service.ImageService
	serviceUser  *service.UserService
	jwt          *middlewares.AuthMiddleware
}

func NewController(app *fiber.App, serviceImage *service.ImageService, serviceUser *service.UserService, jwtSecret string) *Controller {
	return &Controller{
		app:          app,
		serviceImage: serviceImage,
		serviceUser:  serviceUser,
		jwt:          middlewares.NewAuthMiddleware(jwtSecret),
	}
}

func (c *Controller) SetupRoutes() {
	c.getImage()
	c.uploadImage()
	c.deleteImage()
	c.loginUser()
	c.registerUser()
}

func (c *Controller) getImage() {
	c.app.Get("/getImage/:id", c.jwt.Handler(), func(ctx *fiber.Ctx) error {
		token := ctx.Locals("user").(*jtoken.Token)
		dtoUserJWT, errJWT := middlewares.GetJWTClaims(token)
		if errJWT != nil {
			return ctx.Status(errJWT.Status).JSON("Error al validar el usuario")
		}

		_, errUser := c.serviceUser.FindAndCheckJWT(dtoUserJWT)
		if errUser != nil {
			return ctx.Status(errUser.Status).JSON(errUser)
		}

		dtoFindImage := &dto.DTOImage{
			IdImage: ctx.Params("id"),
			Owner:   dtoUserJWT.Username,
		}

		image, err := c.serviceImage.Find(dtoFindImage)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		return ctx.Status(fiber.StatusOK).JSON(image)
	})
}

func (c *Controller) uploadImage() {
	c.app.Post("/uploadImage", c.jwt.Handler(), func(ctx *fiber.Ctx) error {
		token := ctx.Locals("user").(*jtoken.Token)
		dtoUserJWT, errJWT := middlewares.GetJWTClaims(token)
		if errJWT != nil {
			return ctx.Status(errJWT.Status).JSON("Error al validar el usuario")
		}

		_, errUser := c.serviceUser.FindAndCheckJWT(dtoUserJWT)
		if errUser != nil {
			return ctx.Status(errUser.Status).JSON(errUser)
		}

		fileInput, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(fiber.StatusNotFound).JSON(exception.NewApiException(fiber.StatusNotFound, "Error al obtener la imagen del formulario"))
		}

		dtoInsertImage, errFile := handler.ProcessImageFile(fileInput, dtoUserJWT.Username)
		if errFile != nil {
			return ctx.Status(errFile.Status).JSON(errFile)
		}

		dto, errInsert := c.serviceImage.Insert(dtoInsertImage)
		if errInsert != nil {
			return ctx.Status(errInsert.Status).JSON(errInsert)
		}

		return ctx.Status(fiber.StatusOK).JSON(dto)
	})
}

func (c *Controller) deleteImage() {
	c.app.Delete("/deleteImage/:id", c.jwt.Handler(), func(ctx *fiber.Ctx) error {
		token := ctx.Locals("user").(*jtoken.Token)
		dtoUserJWT, errJWT := middlewares.GetJWTClaims(token)
		if errJWT != nil {
			return ctx.Status(errJWT.Status).JSON("Error al validar el usuario")
		}

		_, errUser := c.serviceUser.FindAndCheckJWT(dtoUserJWT)
		if errUser != nil {
			return ctx.Status(errUser.Status).JSON(errUser)
		}

		id := ctx.Params("id")
		image, err := c.serviceImage.Delete(id)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		return ctx.Status(fiber.StatusOK).JSON(image)
	})
}

func (c *Controller) loginUser() {
	c.app.Post("/login", func(ctx *fiber.Ctx) error {
		dtoLoginRequest := new(dto.DTOUser)
		err := ctx.BodyParser(dtoLoginRequest)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erroneo"))
		}

		dtoUser, errFind := c.serviceUser.Find(dtoLoginRequest)
		if errFind != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(errFind)
		}

		day := time.Hour * 24

		// Create the JWT claims, which includes the user ID and expiry time
		claims := jtoken.MapClaims{
			"username": dtoUser.Username,
			"email":    dtoUser.Email,
			"name":     dtoUser.Firstname,
			"exp":      time.Now().Add(day * 1).Unix(),
		}

		// Create token
		token := jtoken.NewWithClaims(jtoken.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte(c.jwt.GetSecret()))
		if err != nil {
			ctx.Status(fiber.StatusInternalServerError).JSON(exception.NewApiException(fiber.StatusInternalServerError, "Error al generar el token"))
		}
		// Return the token
		return ctx.Status(fiber.StatusOK).JSON(dto.DTOLoginResponse{
			Token: t,
		})
	})
}

func (c *Controller) registerUser() {
	c.app.Post("/register", func(ctx *fiber.Ctx) error {
		dtoRegisterRequest := new(dto.DTOUser)
		err := ctx.BodyParser(dtoRegisterRequest)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(exception.NewApiException(fiber.StatusBadRequest, "El JSON enviado en la petición es erroneo"))
		}

		errHandler := handler.ProcessUser(dtoRegisterRequest)
		if errHandler != nil {
			return ctx.Status(errHandler.Status).JSON(errHandler)
		}

		dtoUser, errInsert := c.serviceUser.Insert(dtoRegisterRequest)
		if errInsert != nil {
			return ctx.Status(errInsert.Status).JSON(errInsert)
		}

		dto := dto.DTORegisterResponse{
			Username:  dtoUser.Username,
			Firstname: dtoUser.Firstname,
			Message:   "Se ha creado el usuario correctamente",
		}
		return ctx.Status(fiber.StatusCreated).JSON(dto)
	})
}
