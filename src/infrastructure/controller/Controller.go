package controller

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/config"
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
	serviceImage *service.ServiceImage
	serviceUser  *service.ServiceUser
	jwt          func(*fiber.Ctx) error
}

func NewController(app *fiber.App, serviceImage *service.ServiceImage, serviceUser *service.ServiceUser) *Controller {
	return &Controller{
		app:          app,
		serviceImage: serviceImage,
		serviceUser:  serviceUser,
		jwt:          middlewares.NewAuthMiddleware(config.Secret),
	}
}

func (c *Controller) SetupRoutes() {
	c.getImage()
	c.uploadImage()
	c.deleteImage()
	c.login()
}

func (c *Controller) getImage() {
	c.app.Get("/getImage/:id", c.jwt, func(ctx *fiber.Ctx) error {
		user := ctx.Locals("user").(*jtoken.Token)
		claims := user.Claims.(jtoken.MapClaims)
		//TODO mirar usuario y validar
		username := claims["username"].(string)
		var _ = username

		id := ctx.Params("id")
		image, err := c.serviceImage.Find(id)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		return ctx.Status(200).JSON(image)
	})
}

func (c *Controller) uploadImage() {
	c.app.Post("/uploadImage", func(ctx *fiber.Ctx) error {
		fileInput, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(404).JSON(exception.NewApiException(404, "Error al obtener la imagen del formulario"))
		}

		processedImage, errFile := handler.ProcessImageFile(fileInput)
		if errFile != nil {
			return ctx.Status(errFile.Status).JSON(err)
		}

		dto, errInsert := c.serviceImage.Insert(processedImage)
		if errInsert != nil {
			return ctx.Status(errInsert.Status).JSON(err)
		}

		return ctx.Status(200).JSON(dto)
	})
}

func (c *Controller) deleteImage() {
	c.app.Delete("/deleteImage/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		image, err := c.serviceImage.Delete(id)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		return ctx.Status(200).JSON(image)
	})
}

func (c *Controller) login() {

	loginRequest := new(dto.DTOLoginRequest)
	c.app.Post("/login", func(ctx *fiber.Ctx) error {

		if err := ctx.BodyParser(loginRequest); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Find the user by credentials
		dtoUser, errFind := c.serviceUser.Find(loginRequest.Username, loginRequest.Password)
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
		t, err := token.SignedString([]byte(config.Secret))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Return the token
		return ctx.JSON(dto.DTOLoginResponse{
			Token: t,
		})
	})
}
