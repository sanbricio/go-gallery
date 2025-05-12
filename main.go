package main

import (
	"fmt"
	"go-gallery/src/commons/configurator"
	"go-gallery/src/infrastructure/auth"
	imageController "go-gallery/src/infrastructure/controller/image"
	swaggerController "go-gallery/src/infrastructure/controller/swagger"
	userController "go-gallery/src/infrastructure/controller/user"
	userMiddleware "go-gallery/src/infrastructure/controller/user/middlewares"
	log "go-gallery/src/infrastructure/logger"
	"runtime/debug"

	codeGeneratorService "go-gallery/src/service/codeGenerator"
	emailService "go-gallery/src/service/email"
	imageService "go-gallery/src/service/image"
	userService "go-gallery/src/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

var logger log.Logger

// @title						GoGallery
// @version					v1.0.0
// @description				API for managing photo uploads, with authentication
// @contact.name				Support GoGallery
// @contact.email				gogalleryteam@gmail.com
// @BasePath					/api
// @securityDefinitions.apiKey	CookieAuth
// @in							header
// @name						Cookie
func main() {
	// Initialize the Fiber application
	app := fiber.New()

	// Load configuration and dependency container
	configuration, dependencyContainer := configurator.LoadConfiguration()

	logger = log.Instance()

	// Initialize the EmailSender, User, and Image services
	logger.Info("Initializing EmailSender service...")
	emailSenderService := emailService.NewEmailSenderService(dependencyContainer.GetEmailSenderRepository())

	logger.Info("Initializing User service...")
	userService := userService.NewUserService(dependencyContainer.GetUserRepository())

	logger.Info("Initializing Code Generator service...")
	codeGeneratorService := codeGeneratorService.NewCodeGeneratorService(dependencyContainer.GetCodeGeneratorRepository())

	logger.Info("Initializing Image service...")
	imageService := imageService.NewImageService(dependencyContainer.GetImageRepository(), dependencyContainer.GetThumbnailImageRepository())

	logger.Info("Starting controller configuration...")

	logger.Info("Setting up CORS middleware...")
	// Middleware to allow CORS
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))

	// Add Recover middleware to handle panics and provide custom stack trace handler
	app.Use(recover.New(recover.Config{
		EnableStackTrace: true,
		StackTraceHandler: func(c *fiber.Ctx, e any) {
			// Log the panic with stack trace
			stackTrace := string(debug.Stack())
			logger.Panic(fmt.Sprintf("Recovered from panic: %v\nStack Trace: %s", e, stackTrace))

			// Respond with a generic error message
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error. Please try again later.",
			})
		},
	}))

	// Configure routes for Swagger documentation
	logger.Info("Setting up Swagger documentation routes...")
	docsController := swaggerController.NewSwaggerController(configuration.GetSwaggerConfiguration())
	docsGroup := app.Group("/api/docs")
	docsController.SetUpRoutes(docsGroup)

	// Initialize JWT authentication middleware
	logger.Info("Initializing JWT middleware...")
	tokenManager := auth.NewJWTTokenManager(configuration.GetJWTSecret())
	jwtMiddleware := userMiddleware.NewJWTMiddleware(tokenManager, userService)

	// Configure user authentication routes
	logger.Info("Setting up user authentication routes...")
	authController := userController.NewAuthController(userService, emailSenderService, imageService, codeGeneratorService, jwtMiddleware)
	authGroup := app.Group("/api/auth")
	authController.SetUpRoutes(authGroup)

	// Configure image routes protected by JWT
	logger.Info("Setting up image routes protected by JWT...")
	imageController := imageController.NewImageController(imageService, userService)
	imageGroup := app.Group("/api/image")
	imageGroup.Use(jwtMiddleware.Handler())
	imageController.SetUpRoutes(imageGroup)

	// Start the server and listen on the configured port
	port := configuration.GetPort()
	logger.Info("Starting the server on port: " + port + "...")
	err := app.Listen(":" + port)
	if err != nil {
		logger.Panic("Failed to start the server: " + err.Error())
	}
}
