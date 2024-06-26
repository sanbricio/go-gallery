package main

import (
	"api-upload-photos/src/infrastructure/controller"
	infrastructure "api-upload-photos/src/infrastructure/repository"
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
)

func main() {

	app := fiber.New()
	
	repository := &infrastructure.RepositoryMemory{}
	service := service.NewService(repository)
	controller := controller.NewController(app, service)
	controller.SetupRoutes()
	
	app.Listen(":3000")
}
