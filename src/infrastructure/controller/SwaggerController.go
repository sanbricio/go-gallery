package controller

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

type SwaggerController struct {
	swaggerConfiguration swagger.Config
}

func NewSwaggerController(config swagger.Config) *SwaggerController {
	return &SwaggerController{
		swaggerConfiguration: config,
	}
}
func (c *SwaggerController) SetUpRoutes(router fiber.Router) {
	router.Get("/definition/swagger.json", c.ServeJSON)
	router.Get("/definition/swagger.yml", c.ServeYAML)
	router.Get("/*", c.ServeSwagger)
}

func (c *SwaggerController) ServeSwagger(ctx *fiber.Ctx) error {
	return swagger.New(c.swaggerConfiguration)(ctx)
}

func (c *SwaggerController) ServeYAML(ctx *fiber.Ctx) error {
	fileContent, err := os.ReadFile("./docs/swagger.yaml")
	if err != nil {
		log.Println("Error abriendo el archivo swagger.yaml:", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error al cargar la documentación YAML")
	}

	return ctx.Send(fileContent)
}

func (c *SwaggerController) ServeJSON(ctx *fiber.Ctx) error {
	fileContent, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Println("Error abriendo el archivo swagger.json:", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error al cargar la documentación JSON")
	}

	return ctx.Send(fileContent)
}
