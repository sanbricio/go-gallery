package controller

import (
	"api-upload-photos/src/commons/exception"
	handler "api-upload-photos/src/infrastructure"
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	app     *fiber.App
	service *service.Service
}

func NewController(app *fiber.App, service *service.Service) *Controller {
	return &Controller{
		app:     app,
		service: service,
	}
}

func (c *Controller) SetupRoutes() {
	c.getIndexPage()
	c.getImage()
	c.uploadImage()
	c.deleteImage()
}

func (c *Controller) getIndexPage() {
	c.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})
}

func (c *Controller) getImage() {
	c.app.Get("/getImage/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		image, err := c.service.Find(id)
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

		dto, errInsert := c.service.Insert(processedImage)
		if errInsert != nil {
			return ctx.Status(errInsert.Status).JSON(err)
		}

		return ctx.Status(200).JSON(dto)
	})
}

func (c *Controller) deleteImage() {
	c.app.Delete("/deleteImage/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		image, err := c.service.Delete(id)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		return ctx.Status(200).JSON(image)
	})
}
