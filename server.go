package main

import (
	"api-upload-photos/src/infrastructure/controller"
	infrastructure "api-upload-photos/src/infrastructure/repository"
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	})
	
	repository := &infrastructure.RepositoryMemory{}
	service := service.NewService(repository)
	controller := controller.NewController(app, service)
	controller.SetupRoutes()
	
	app.Listen(":3000")
}
