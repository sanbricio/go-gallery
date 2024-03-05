package main

import (
	"api-upload-photos/src/infrastructure"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{})
	})

	app.Post("/uploadImage", func(c *fiber.Ctx) error {
		fileInput, err := c.FormFile("file")
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error al obtener la imagen del formulario")
		}

		response, err := infrastructure.AddImageToDataBase(fileInput)

		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Error: "+err.Error())
		}
		return c.Render("status", fiber.Map{
			"Response": response,
		})
	})

	app.Get("/getImage/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		image, err := infrastructure.GetImageByID(id)
		if err != nil {
			return fiber.NewError(fiber.StatusNotFound, "Imagen no encontrada")
		}

		return c.Render("showImage", fiber.Map{
			"Response": image,
		})
	})

	log.Fatal(app.Listen(":3000"))
}
