package main

import (
	"api-upload-photos/src/commons/configurator"
	"api-upload-photos/src/infrastructure/controller"
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))
	// Cargamos la configuración y el contenedor de dependencias
	configuration, dependencyContainer := configurator.LoadConfiguration()

	serviceImage := service.NewImageService(dependencyContainer.GetImageRepository())
	serviceUser := service.NewUserService(dependencyContainer.GetUserRepository())
	controller := controller.NewController(app, serviceImage, serviceUser, configuration.GetJWTSecret())
	controller.SetupRoutes()

	app.Listen(":" + configuration.GetPort())
}
