package main

import (
	"go-gallery/src/commons/configurator"
	imageController "go-gallery/src/infrastructure/controller/image"
	swaggerController "go-gallery/src/infrastructure/controller/swagger"
	userController "go-gallery/src/infrastructure/controller/user"
	userMiddleware "go-gallery/src/infrastructure/controller/user/middlewares"

	emailService "go-gallery/src/service/email"
	imageService "go-gallery/src/service/image"
	userService "go-gallery/src/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

//	@title						GoGallery
//	@version					v1.0.0
//	@description				API para la gestión de subida de fotos, con una autentificación
//	@contact.name				Support GoGallery
//	@contact.email				gogalleryteam@gmail.com
//	@BasePath					/api
//	@securityDefinitions.apiKey	CookieAuth
//	@in							header
//	@name						Cookie
func main() {
	// Inicializamos la aplicación Fiber
	app := fiber.New()

	// Cargamos la configuración y el contenedor de dependencias
	configuration, dependencyContainer := configurator.LoadConfiguration()

	// Inicializamos los servicios de emailSender, usuario e imagen
	emailSenderService := emailService.NewEmailSenderService(dependencyContainer.GetEmailSenderRepository())
	userService := userService.NewUserService(dependencyContainer.GetUserRepository())
	imageService := imageService.NewImageService(dependencyContainer.GetImageRepository(), dependencyContainer.GetThumbnailImageRepository())

	// Middleware para permitir CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Configuración de rutas para la documentación de Swagger
	docsController := swaggerController.NewSwaggerController(configuration.GetSwaggerConfiguration())
	docsGroup := app.Group("/api/docs")
	docsController.SetUpRoutes(docsGroup)

	// Inicializamos el middleware de autenticación JWT
	jwtMiddleware := userMiddleware.NewJWTMiddleware(configuration.GetJWTSecret())

	// Configuración de rutas de autenticación de usuario
	authController := userController.NewAuthController(userService, emailSenderService, jwtMiddleware)
	authGroup := app.Group("/api/auth")
	authController.SetUpRoutes(authGroup)

	// Configuración de rutas de imágenes con autenticación JWT
	imageController := imageController.NewImageController(imageService, userService)
	imageGroup := app.Group("/api/image")
	imageGroup.Use(jwtMiddleware.Handler())
	imageController.SetUpRoutes(imageGroup)

	// Iniciamos el servidor escuchando en el puerto configurado
	app.Listen(":" + configuration.GetPort())
}
