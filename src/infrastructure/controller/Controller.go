package controller

import (
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	app * fiber.App
	service *service.Service
}

func NewController(app *fiber.App, service *service.Service) *Controller {
	return &Controller{
		app:  app,
		service: service,
	}
}

func (c *Controller) SetupRoutes() {
	c.getIndexPage()
	c.setupUploadImage()
	c.setupGetImage()
}

func (c *Controller) getIndexPage() {
	c.app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})
}

func (controller *Controller) setupUploadImage() {
	controller.app.Post("/uploadImage", func(c *fiber.Ctx) error {
		fileInput, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error al obtener la imagen del formulario")
		}

		response, err := controller.service.Insert(fileInput)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error: "+err.Error())
		}
		return c.Render("status", fiber.Map{
			"Response": response,
		})
	})
}

func (controller *Controller) setupGetImage() {
	controller.app.Get("/getImage/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		image, err := controller.service.Find(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Imagen no encontrada")
		}

		return c.Render("showImage", fiber.Map{
			"Response": image,
		})
	})
}
