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

	// Inicializamos los servicios de imagen y usuario
	imageService := service.NewImageService(dependencyContainer.GetImageRepository())
	userService := service.NewUserService(dependencyContainer.GetUserRepository())

	// Middleware para permitir CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Inicializamos el middleware de autenticación JWT
	authMiddleware := middlewares.NewAuthMiddleware(configuration.GetJWTSecret())

	// Middleware para validar el token JWT en todas las rutas excepto en /auth/login y /auth/register
	app.Use(func(ctx *fiber.Ctx) error {
		// Excluimos las rutas de login y register de la validación JWT
		if ctx.Path() == "/auth/login" || ctx.Path() == "/auth/register" {
			return ctx.Next()
		}
		return authMiddleware.Handler()(ctx)
	})

	// Configuración de rutas de autenticación
	authController := controller.NewAuthController(userService, configuration.GetJWTSecret())
	authGroup := app.Group("/auth")
	authController.SetUpRoutes(authGroup)

	// Configuración de rutas de imágenes
	imageController := controller.NewImageController(imageService, userService)
	imageGroup := app.Group("/image")
	imageController.SetUpRoutes(imageGroup)

	// Iniciamos el servidor escuchando en el puerto configurado
	app.Listen(":" + configuration.GetPort())
}
