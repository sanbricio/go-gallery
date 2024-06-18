package controller

import (
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
	c.setupUploadImage()
	c.setupGetImage()
}

func (c *Controller) getIndexPage() {
	c.app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.Render("index", fiber.Map{})
	})
}

func (c *Controller) setupUploadImage() {
	c.app.Post("/uploadImage", func(ctx *fiber.Ctx) error {
		fileInput, err := ctx.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error al obtener la imagen del formulario")
		}

		response, errInsert := c.service.Insert(fileInput)
		if errInsert != nil {
			return fiber.NewError(errInsert.GetStatus(), errInsert.GetMessage())
		}

		dto := dto.FromResponse(response)

		return ctx.Render("status", fiber.Map{
			"Response": dto,
		})
	})
}

func (c *Controller) setupGetImage() {
	c.app.Get("/getImage/:id", func(ctx *fiber.Ctx) error {
		id := ctx.Params("id")
		image, err := c.service.Find(id)
		if err != nil {
			return ctx.Render("showImage", fiber.Map{
				"Status": err.GetStatus(),
				"Error":  err.GetMessage(),
			})
		}

		return ctx.Render("showImage", fiber.Map{
			"Response": image,
		})
	})
}
