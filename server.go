package main

import (
	"api-upload-photos/src/config"
	"api-upload-photos/src/infrastructure/controller"
	repositoryImage "api-upload-photos/src/infrastructure/repository/image"
	repositoryUser "api-upload-photos/src/infrastructure/repository/user"
	"api-upload-photos/src/service"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error cargando el archivo .env: %v", err)
	}
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOriginsFunc: func(origin string) bool {
			return true
		},
	}))
	//TODO ¿llevar al configurator?
	config.Secret = os.Getenv("SECRET_KEY")
	mongoURL := os.Getenv("LOCAL_MONGODB_URL")
	databaseName := os.Getenv("MONGODB_DATABASE")

	repositoryImage, errRepoImage := repositoryImage.NewImageMongoDBRepository(mongoURL, databaseName)
	if errRepoImage != nil {
		log.Fatalf("[ERROR] %s\n StackTrace:\n%s", errRepoImage.Message, errRepoImage.StackTrace)
	}

	repositoryUser, errRepoUser := repositoryUser.NewUserMongoDBRepository(mongoURL, databaseName)
	if errRepoUser != nil {
		log.Fatalf("[ERROR] %s\n StackTrace:\n%s", errRepoUser.Message, errRepoUser.StackTrace)
	}

	serviceImage := service.NewServiceImage(repositoryImage)
	serviceUser := service.NewServiceUser(repositoryUser)
	controller := controller.NewController(app, serviceImage, serviceUser)
	controller.SetupRoutes()

	app.Listen(":3000")
}
