package controller

import (
	"api-upload-photos/src/commons/exception"
	"api-upload-photos/src/domain/dto"
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
			return ctx.Render("showImage", fiber.Map{
				"Status": err.Status,
				"Error":  err.Message,
			})
		}

		return ctx.Render("showImage", fiber.Map{
			"Response": image,
		})
	})
}

// TODO Mirar para futuro desacople como se haria el ctx.FormFile("file") si no tengo el front en el mismo proyecto
func (c *Controller) uploadImage() {
	c.app.Post("/uploadImage", func(ctx *fiber.Ctx) error {
		fileInput, err := ctx.FormFile("file")
		if err != nil {
			return ctx.Status(404).JSON(exception.NewApiException(404, "Error al obtener la imagen del formulario"))
		}

		response, errInsert := c.service.Insert(fileInput)
		if errInsert != nil {
			return ctx.Status(errInsert.Status).JSON(err)
		}

		dto := dto.FromResponse(response)

		return ctx.Render("status", fiber.Map{
			"Response": dto,
		})
	})
}

func (c *Controller) deleteImage() {
	c.app.Delete("/deleteImage/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		image, err := c.service.Delete(id)
		if err != nil {
			return ctx.Status(err.Status).JSON(err)
		}

		return ctx.JSON(image)
	})
}
