package main

import (
	"api-upload-photos/src/commons/configurator"
	"api-upload-photos/src/infrastructure/controller"
	"api-upload-photos/src/infrastructure/controller/middlewares"
	"api-upload-photos/src/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Inicializamos la aplicación Fiber
	app := fiber.New()

	// Cargamos la configuración y el contenedor de dependencias
	configuration, dependencyContainer := configurator.LoadConfiguration()

	// Inicializamos los servicios de emailSender, usuario e imagen
	emailSenderService := service.NewEmailSenderService(dependencyContainer.GetEmailSenderRepository())
	userService := service.NewUserService(dependencyContainer.GetUserRepository())
	imageService := service.NewImageService(dependencyContainer.GetImageRepository())

	// Middleware para permitir CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Inicializamos el middleware de autenticación JWT
	jwtMiddleware := middlewares.NewJWTMiddleware(configuration.GetJWTSecret())

	// Configuración de rutas de autenticación de usuario
	authController := controller.NewAuthController(userService, emailSenderService, jwtMiddleware)
	authGroup := app.Group("/auth")
	authController.SetUpRoutes(authGroup)

	// Configuración de rutas de imágenes con autenticación JWT
	imageController := controller.NewImageController(imageService, userService)
	imageGroup := app.Group("/image")
	imageGroup.Use(jwtMiddleware.Handler())
	imageController.SetUpRoutes(imageGroup)

	// Iniciamos el servidor escuchando en el puerto configurado
	app.Listen(":" + configuration.GetPort())
}
