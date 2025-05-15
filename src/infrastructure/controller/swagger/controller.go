package swaggerController

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

// @Summary		Obtiene la documentación de la API en formato YAML
// @Description	Retorna la definición de la API(OpenAPI) en formato YAML
// @Tags			docs
// @Produce		plain
// @Success		200	"Archivo YAML cargado correctamente"
// @Failure		500	"Error al cargar el archivo YAML"
// @Router			/docs/definition/swagger.yml [get]
func (c *SwaggerController) ServeYAML(ctx *fiber.Ctx) error {
	fileContent, err := os.ReadFile("./docs/swagger.yaml")
	if err != nil {
		log.Println("Error abriendo el archivo swagger.yaml:", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error al cargar la documentación YAML")
	}

	return ctx.Send(fileContent)
}

// @Summary		Obtiene la documentación de la API en formato JSON
// @Description	Retorna la definición de la API(OpenAPI) en formato JSON
// @Tags			docs
// @Produce		json
// @Success		200	"Archivo JSON cargado correctamente"
// @Failure		500	"Error al cargar el archivo JSON"
// @Router			/docs/definition/swagger.json [get]
func (c *SwaggerController) ServeJSON(ctx *fiber.Ctx) error {
	fileContent, err := os.ReadFile("./docs/swagger.json")
	if err != nil {
		log.Println("Error abriendo el archivo swagger.json:", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error al cargar la documentación JSON")
	}

	return ctx.Send(fileContent)
}
